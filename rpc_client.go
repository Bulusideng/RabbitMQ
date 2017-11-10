package main

import (
	"amqp"
	"encoding/json"

	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func rpc_client() {
	rand.Seed(time.Now().UTC().UnixNano())
	conn, err := amqp.Dial(host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"exchange_1", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	responseQ, err := ch.QueueDeclare(
		"rsp_queue", // name
		false,       // durable
		false,       // delete when usused
		true,        // exclusive
		false,       // noWait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		responseQ.Name,  // queue name
		"svc1_response", // routing key
		"exchange_1",    // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	responses, err := ch.Consume(
		responseQ.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	req := &Request{
		"test",
	}

	rsp := &Response{}
	for {
		//n := randInt(1, 30)

		time.Sleep(time.Second * 1)

		req.CorrelationId = randomString(32)
		req.MachineNum += 1
		jmsg, err := json.Marshal(req)

		err = ch.Publish(
			"exchange_svc1", // exchange
			"",              // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: req.CorrelationId,
				ReplyTo:       "svc1_response", //respone rounting key
				Body:          []byte(jmsg),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf("[CLI]->%s: %+v", req.CorrelationId, *req)

		for rcvd := range responses {
			if req.CorrelationId == rcvd.CorrelationId {
				err = json.Unmarshal(rcvd.Body, rsp)
				log.Printf("[CLI]<-%s: %+v", rsp.CorrelationId, *rsp)
				log.Println()
				break
			} else {
				failOnError(err, "Unexpected corrid")
			}
		}
		failOnError(err, "Failed to handle RPC request")
	}

}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (res int, err error) {

	return
}

func bodyFrom1(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	failOnError(err, "Failed to convert arg to integer")
	return n
}
