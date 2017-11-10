package main

import (
	"amqp"
)

type ServiceInterface struct {
	name       string
	deliveries <-chan amqp.Delivery
}

func (this *ServiceInterface) Process() {

	return &Response{}
}

func (this *ServiceInterface) Publish() {

}
