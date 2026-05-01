package queues

import (
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

func RunBackgroundConsumer() (map[string]<-chan amqp091.Delivery, error) {

	channelsResult := make(map[string]<-chan amqp091.Delivery, 0)
	for queueKey, actualQueue := range Queues {

		if !actualQueue.IsConsumer {
			continue
		}
		msgs, err := QueueMainCh.Consume(
			actualQueue.QueueObject.Name, // queue
			"",                           // consumer
			false,                        // auto-ack
			false,                        // exclusive
			false,                        // no-local
			false,                        // no-wait
			nil,                          // args
		)
		if err != nil {
			return nil, errors.New("Failed to register a consumer with name " + actualQueue.QueueObject.Name + " because: " + err.Error())
		}
		channelsResult[queueKey] = msgs
	}
	return channelsResult, nil
}
