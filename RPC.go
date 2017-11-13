package main

import (
	"encoding/json"

	"rabitmq/amqp"
)

type RPC struct {
	host       string
	exchange   string
	conn       *amqp.Connection
	ch         *amqp.Channel
	queues     []*RPCQueue
	pubChan    chan *Msg
	MaxPubSize int
}

func NewRPC(host, exchange string) *RPC {
	return &RPC{
		host:     host,
		exchange: EXPrefix + exchange,
	}
}

func (this *RPC) AddQueue(qname, iname string) { //Add queue for interface
	q := &RPCQueue{
		name:    qname,
		bindKey: iname, //The binding key is the interface name
	}
	this.queues = append(this.queues, q)
}

func (this *RPC) Send() {
	jmsg, err := json.Marshal(response)
	failOnError(err, "Failed to convert body to integer")

	err = ch.Publish(
		this.exchange, // exchange
		d.ReplyTo,     // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          jmsg,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("[SVR]->%s: %+v", d.CorrelationId, *response)
}

func (this *RPC) Destroy() {
	if this.ch != nil {
		this.ch.Close()
	}
	if this.conn != nil {
		this.conn.Close()
	}
}

func (this *RPC) Init() (err error) {
	this.conn, err = amqp.Dial(host)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		this.Destroy()
		return
	}

	this.ch, err = this.conn.Channel()
	if err != nil {
		this.Destroy()
		failOnError(err, "Failed to open a channel")
		return err
	}
	err = this.ch.ExchangeDeclare(
		this.exchange, // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		this.Destroy()
		failOnError(err, "Failed to declare an exchange")
		return err
	}

	for _, q := range this.queues {
		q.queue, err = this.ch.QueueDeclare(
			q.name, // name
			false,  // durable
			false,  // delete when usused
			true,   // exclusive
			false,  // noWait
			nil,    // arguments
		)

		if err != nil {
			this.Destroy()
			failOnError(err, "Failed to declare a queue")
			return err
		}

		err = this.ch.QueueBind(
			q.name,
			q.bindKey,
			this.exchange,
			false,
			nil)

		if err != nil {
			this.Destroy()
			failOnError(err, "Failed to bind a queue")
			return err
		}

		q.Deliveries, err = this.ch.Consume(
			q.name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			this.Destroy()
			failOnError(err, "Failed to register a consumer")
			return err
		}
	}
	return nil
}

type Msg struct {
	routKey string
	corid   string
	replyTo string
	body    interface{}
}

func (this *RPC) doPub() {
	for {
		if jmsg, ok := <-this.pubChan; ok {
			msg, err := json.Marshal(jmsg)

			err = this.ch.Publish(
				this.exchange,
				jmsg.routKey,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: jmsg.corid,
					ReplyTo:       jmsg.replyTo,
					Body:          []byte(msg),
				})
			if err != nil {
				break
			}
		} else {
			break
		}
	}
}

func (this *RPC) Publish(msg *Msg) bool {
	if len(this.pubChan) >= this.MaxPubSize {
		return false
	}

	this.pubChan <- msg
	return true

}
