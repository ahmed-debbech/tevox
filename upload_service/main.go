package main

import (
	"log"
)

func main() {
	log.Println("Upload Service is up & running!")

	err := ConnectAndPrepare()
	if err != nil {
		log.Fatal(err)
	}
	var consumer Consumer = ConsumerJob{}
	HookConsumers(consumer)
	select {}
}
