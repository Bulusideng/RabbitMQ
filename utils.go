package main

import (
	"log"
)

type Request struct {
	Cmd           string
	Interface     string
	CorrelationId string
	GUid          string
	ReplyTo       string
	Body          []interface{}
}

type Response struct {
	Cmd           string
	Interface     string
	CorrelationId string
	Body          interface{}
}

var (
	host      = "amqp://guest:guest@localhost:5672/"
	XN_FANOUT = "exchange_fanout"
	XN_DIRECT = "exchange_topic"
	XN_TOPIC  = "exchange_topic"
	XN_HEADER = "exchange_header"
)

var (
	XT_FANOUT = "fanout"  //Broadcast
	XT_DIRECT = "direct"  //Match key
	XT_TOPIC  = "topic"   //dot split key with * and # match
	XT_HEADER = "headers" //match headers
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
