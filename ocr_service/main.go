package main

import (
	"log"

	"github.com/ahmed-debbech/tevox/ocr_service/core"
	"github.com/ahmed-debbech/tevox/ocr_service/queues"
)

func main() {
	log.Println("OCR SERVICE started")

	if err := queues.DefineAllQueues(); err != nil {
		log.Fatal(err)
	}

	mapOfConsumerCh, err := queues.RunBackgroundConsumer()
	if err != nil {
		log.Fatal()
	}
	core.HookConsumers(mapOfConsumerCh)
	select {}
}
