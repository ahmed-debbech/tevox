package queues

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(queueKey string, content []byte, headers map[string]interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := QueueMainCh.PublishWithContext(ctx,
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
	QueueMainCh.Nack(deliveryTag, false, false)
}

func Ack(deliveryTag uint64) {
	QueueMainCh.Ack(deliveryTag, false)
}
