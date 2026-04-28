package core

import (
	"log"
	"strconv"

	"github.com/ahmed-debbech/tevox/ocr_service/config"
	"github.com/ahmed-debbech/tevox/ocr_service/queues"
	"github.com/rabbitmq/amqp091-go"
)

func HookConsumers(consumers map[string]<-chan amqp091.Delivery) {

	go ScanImageEventConsumer(consumers["A_QUEUE"], "A_QUEUE")

}

func ScanImageEventConsumer(feed <-chan amqp091.Delivery, queueName string) {

	qConf := queues.GetQueueConfig(queueName)

	for msg := range feed {

		log.SetPrefix("DELIVERY[" + strconv.Itoa(int(msg.DeliveryTag)) + "] ")
		log.Println("Message received from queue", msg.Body)

		if val, ok := msg.Headers["x-death"]; ok {

			retry_count := queues.GetRetryNumber(val, qConf.Name)
			log.Println("checking if message retry number exceeded max, current:", retry_count, "max:", config.RetryCount)
			if retry_count > config.RetryCount {
				msg.Ack(false)
				log.Println("message reached maximum retry counts - discarding...")
				log.Println("publishing to DLQ")
				err := queues.Publish(qConf.Dlq.Name, msg.Body, msg.Headers)
				if err != nil {
					log.Println("could not publish failed message to DLQ")
				}
				continue
			}
		}

		if err := ProcessScanImageEvent(msg.Body); err != nil {
			log.Println(err)
			msg.Nack(false, false)
		} else {
			msg.Ack(false)
		}
		log.SetPrefix("")
	}

}
