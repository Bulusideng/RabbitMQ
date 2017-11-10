package main

import (
	"amqp"
	"encoding/json"
	"log"
	"strconv"
)

func rpc_server() {
	conn, err := amqp.Dial(host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"exchange_svc1", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	requestQ, err := ch.QueueDeclare(
		"req_queue", // name
		false,       // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		requestQ.Name,   // queue name
		"",              // routing key
		"exchange_svc1", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	requests, err := ch.Consume(
		requestQ.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	request := &Request{}
	response := &Response{}
	go func() {
		for d := range requests {
			if false {
				_, err := strconv.Atoi(string(d.Body))
				failOnError(err, "Failed to convert body to integer")
			} else {
				err := json.Unmarshal(d.Body, request)
				failOnError(err, "Failed to convert body to integer")
				log.Printf("[SVR]<-%s: %+v", d.CorrelationId, *request)
			}

			response.Cmd = request.Cmd + "Resp"
			response.Interface = request.Interface
			response.CorrelationId = request.CorrelationId
			response.Body = Status{"Success"}

			jmsg, err := json.Marshal(response)
			failOnError(err, "Failed to convert body to integer")

			err = ch.Publish(
				"exchange_svc1", // exchange
				d.ReplyTo,       // routing key
				false,           // mandatory
				false,           // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          jmsg,
				})
			failOnError(err, "Failed to publish a message")
			log.Printf("[SVR]->%s: %+v", d.CorrelationId, *response)
			d.Ack(false)
		}
	}()

	<-forever
	log.Printf(" [.] exit")
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return n * n //fib(n-1) + fib(n-2)
	}
}
