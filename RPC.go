package main

import (
	"amqp"
	"encoding/json"
)

type RPCQueue struct {
	name    string
	bindKey string
	queue   amqp.Queue
}

type RPC struct {
	host       string
	exchange   string
	conn       *amqp.Connection
	ch         *amqp.Channel
	queues     []*RPCQueue
	pubChan    chan *Msg
	MaxPubSize int
}

func (this *RPC) AddQueue(name, bindKey string) {
	q := &RPCQueue{name: name, bindKey: bindKey}
	this.queues = append(this.queues, q)
}

func (this *RPC) Stop() {
	if this.ch != nil {
		this.ch.Close()
	}
	if this.conn != nil {
		this.conn.Close()
	}
}

func (this *RPC) Start() (err error) {
	this.conn, err = amqp.Dial(host)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		this.Stop()
		return
	}

	this.ch, err = this.conn.Channel()
	if err != nil {
		this.Stop()
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
		this.Stop()
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
			this.Stop()
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
			this.Stop()
			failOnError(err, "Failed to bind a queue")
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
				msg.routKey,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: msg.corid,
					ReplyTo:       msg.replyTo,
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

}

func NewRPC(host, exchange string) *RPC {
	return &RPC{
		host:     host,
		exchange: EXPrefix + exchange,
	}
}
