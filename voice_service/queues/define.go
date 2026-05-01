package queues

import (
	"errors"

	"github.com/ahmed-debbech/tevox/voice_service/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Queues                    = make(map[string]QueueMetadata)
	QueueMainCh *amqp.Channel = nil
)

type QueueMetadata struct {
	QueueObject amqp.Queue
	IsConsumer  bool
}

func DefineAllQueues() error {
	conn, err := amqp.Dial(config.RabbitMqUrl)
	if err != nil {
		return errors.New("Failed to connect to RabbitMQ because," + err.Error())
	}

	QueueMainCh, err = conn.Channel()
	if err != nil {
		return errors.New("Failed to open a channel")
	}

	for queueKey, queueConf := range config.QueuesNames {

		//check if there is DLQ configured for this queue
		if queueConf.Dlq.Name != "" {

			_, err := QueueMainCh.QueueDeclare(
				queueConf.Dlq.Name, // name
				true,               // durability
				false,              // delete when unused
				false,              // exclusive
				false,              // no-wait
				amqp.Table{},
			)
			if err != nil {
				return errors.New("Could not define " + queueConf.Name + " because:" + err.Error())
			}
		}

		//check if there is a RETRY queue configured
		if queueConf.Retry.Name != "" {
			_, err := QueueMainCh.QueueDeclare(
				queueConf.Retry.Name, // name
				true,                 // durability
				false,                // delete when unused
				false,                // exclusive
				false,                // no-wait
				amqp.Table{
					"x-message-ttl":             queueConf.Retry.Ttl,
					"x-dead-letter-exchange":    "",
					"x-dead-letter-routing-key": queueConf.Name,
				},
			)
			if err != nil {
				return errors.New("Could not define " + queueConf.Name + " because:" + err.Error())
			}
		}

		//declare main queue
		q, err := QueueMainCh.QueueDeclare(
			queueConf.Name, // name
			true,           // durability
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			amqp.Table{
				"x-dead-letter-exchange":    "",
				"x-dead-letter-routing-key": queueConf.Retry.Name,
			},
		)
		if err != nil {
			return errors.New("Could not define " + queueConf.Name + " because:" + err.Error())
		}
		qm := QueueMetadata{
			QueueObject: q,
			IsConsumer:  queueConf.Consumer,
		}
		Queues[queueKey] = qm
	}
	return nil
}
