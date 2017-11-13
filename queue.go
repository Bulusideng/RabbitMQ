package main

import (
	"rabitmq/amqp"
)

type RPCQueue struct {
	name       string
	bindKey    string
	queue      amqp.Queue
	Deliveries <-chan amqp.Delivery
}
