package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	MainConnection *amqp.Connection
	MainChannel    *amqp.Channel
	QueuesDefs     []QueueInfo
	x_death_header map[string]int
)

type QueueInfo struct {
	Queue      amqp.Queue
	Consumable bool
}

func ConnectAndPrepare() error {
	conn, err := amqp.Dial(RabbitMqUrl)
	if err != nil {
		log.Println("could not connect to rabbitmq, because", err.Error())
		return errors.New(err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("could not open channel to rabbitmq, because", err)
		return errors.New(err.Error())
	}
	MainChannel = ch

	QueuesDefs = make([]QueueInfo, 0)
	if err := defineInputQueues(); err != nil {
		return err
	}
	return nil
}

func defineInputQueues() error {
	q, err := MainChannel.QueueDeclare(
		"voice_output_q", // name
		true,             // durability
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "voice_output_retry_q",
		},
	)
	if err != nil {
		log.Println("could not define voice_output_q, because", err)
		return errors.New(err.Error())
	}
	QueuesDefs = append(QueuesDefs, QueueInfo{Queue: q, Consumable: true})

	q_retry, err := MainChannel.QueueDeclare(
		"voice_output_retry_q", // name
		true,                   // durability
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		amqp.Table{
			"x-message-ttl":             5000,
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "voice_output_q",
		},
	)
	if err != nil {
		log.Println("could not define voice_output_retry_q, because", err)
		return errors.New(err.Error())
	}
	QueuesDefs = append(QueuesDefs, QueueInfo{Queue: q_retry, Consumable: false})

	q_dlq, err := MainChannel.QueueDeclare(
		"voice_output_dlq_q", // name
		true,                 // durability
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		amqp.Table{},
	)
	if err != nil {
		log.Println("could not define voice_output_dlq_q, because", err)
		return errors.New(err.Error())
	}
	QueuesDefs = append(QueuesDefs, QueueInfo{Queue: q_dlq, Consumable: false})
	return nil
}

func HookConsumers(consumer Consumer) error {

	for _, v := range QueuesDefs {
		if !v.Consumable {
			continue
		}
		msgs, err := MainChannel.Consume(
			v.Queue.Name, // queue
			"",           // consumer
			false,        // auto-ack
			false,        // exclusive
			false,        // no-local
			false,        // no-wait
			nil,          // args
		)
		if err != nil {
			log.Println("could not register queue", v.Queue.Name, "to consume it.")
			return errors.New(err.Error())
		}

		go func() {
			for msg := range msgs {

				log.SetPrefix("DELIVERY[" + strconv.Itoa(int(msg.DeliveryTag)) + "] ")
				log.Println("Message received from queue", v.Queue.Name)

				if x_death_header, ok := msg.Headers["x-death"]; ok {

					retry_count := int64(0)
					for _, val := range x_death_header.([]interface{}) {
						countN, ok1 := val.(amqp091.Table)["count"]
						queueN, ok2 := val.(amqp091.Table)["queue"]
						if ok1 && ok2 {
							if queueN == v.Queue.Name {
								retry_count += countN.(int64)
							}
						}
					}
					log.Println("checking if message retry number exceeded max, current:", retry_count, "max:", RetryCount)
					if retry_count > int64(RetryCount) {
						msg.Ack(false)
						log.Println("message reached maximum retry counts - discarding...")
						dlqName := getDlqNameForQueue(v.Queue.Name)
						if dlqName == "" {
							log.Println("could not find DLQ queue for this queue.")
							continue
						}
						log.Println("publishing to DLQ", dlqName)
						err := Publish(dlqName, msg.Body, msg.Headers)
						if err != nil {
							log.Println("could not publish failed message to DLQ")
						}
						continue
					}
				}

				if err := consumer.Run(msg.Body); err != nil {
					log.Println(err)
					msg.Nack(false, false)
				} else {
					msg.Ack(false)
				}
				log.SetPrefix("")
			}
		}()

		log.Printf("Hooked queue", v.Queue.Name, "! ready for consuming")
	}
	return nil
}

func getDlqNameForQueue(queue string) string {
	for _, v := range QueuesDefs {
		if strings.HasPrefix(v.Queue.Name, queue[:len(queue)-2]) &&
			strings.Contains(v.Queue.Name, "dlq") {
			return v.Queue.Name
		}
	}
	return ""
}

func Publish(queueKey string, content []byte, headers map[string]interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := MainChannel.PublishWithContext(ctx,
		"",       // exchange
		queueKey, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        content,
			Headers:     headers,
		})
	if err != nil {
		return errors.New("Could not publish to queue " + queueKey + " the value " + string(content) + " because " + err.Error())
	}
	return nil
}

func Nack(deliveryTag uint64) {
	MainChannel.Nack(deliveryTag, false, false)
}

func Ack(deliveryTag uint64) {
	MainChannel.Ack(deliveryTag, false)
}
