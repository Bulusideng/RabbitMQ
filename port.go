package main

import (
	"log"
	//	"rabitmq/amqp"
)

type Port interface {
	Start()
}

type ServicePort struct {
	name      string
	inbound   *RPCQueue
	outbound  chan *Request
	MaxOutLen int
}

func (this *ServicePort) Start() {
	go this.handleInbound()
	go this.handleOutbound()
}

func (this *ServicePort) Send(req *Request) bool {
	if len(this.outbound) >= this.MaxOutLen {
		return false
	}
	this.outbound <- req
	return true
}

func (this *ServicePort) handleInbound() {
	for {
		if delivery, ok := <-this.inbound.Deliveries; ok {
			log.Printf("%s receives: %+v", this.name, delivery)
		} else {
			break
		}
	}
}
func (this *ServicePort) handleOutbound() {
	for {
		if req, ok := <-this.outbound; ok {
			log.Printf("%s send: %+v", this.name, req)
		} else {
			break
		}
	}
}
