package queues

import (
	"log"

	"github.com/ahmed-debbech/tevox/voice_service/config"
	"github.com/rabbitmq/amqp091-go"
)

func GetRetryNumber(x_death_header interface{}, queueResultingDeath string) int {
	log.Println(x_death_header)
	c := int64(0)
	for _, val := range x_death_header.([]interface{}) {
		countN, ok1 := val.(amqp091.Table)["count"]
		queueN, ok2 := val.(amqp091.Table)["queue"]
		if ok1 && ok2 {
			if queueN == queueResultingDeath {
				c = c + countN.(int64)
			}
		}
	}
	return int(c)
}

func GetQueueConfig(name string) config.QueueConf {
	return config.QueuesNames[name]
}

func CreateErrorChannel() chan *amqp091.Error {
	return make(chan *amqp091.Error)
}
