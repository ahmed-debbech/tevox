package main

import (
	"log"
	"time"

	"github.com/ahmed-debbech/tevox/ocr_service/core"
	"github.com/ahmed-debbech/tevox/ocr_service/queues"
)

func main() {
	log.Println("OCR SERVICE started")

	i := 0
	for {
		if err := queues.DefineAllQueues(); err != nil {
			if i > 0 {
				log.Println("Trying to restart rabbit mq connection for", i, "times ..")
				i++
				time.Sleep(time.Second * 5)
				continue
			}
			log.Fatal(err)
		}

		mapOfConsumerCh, err := queues.RunBackgroundConsumer()
		if err != nil {
			log.Fatal()
		}

		core.HookConsumers(mapOfConsumerCh)

		chClosed := queues.QueueMainCh.NotifyClose(queues.CreateErrorChannel())

		i = 0
		select {
		case <-chClosed:
		}

		queues.QueueMainCh.Close()
		log.Println("Connection lost to rabbitMq")
		i++
	}
}
