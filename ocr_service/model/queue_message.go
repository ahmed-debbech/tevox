package model

import "github.com/rabbitmq/amqp091-go"

type QueueMessage struct {
	Body    []byte
	Headers amqp091.Table
}
