package main

import (
	"errors"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
)

type QC struct {
	url         string
	exName      string
	exType      string
	routingKeys []string
	consume     bool
	contentType string
	durable     bool

	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(
	url,
	exName,
	exType string,
	routingKeys []string,
	consume bool,
	contentType string,
	durable bool,
) *QC {
	if url == "" {
		panic(errors.New("url required"))
	}
	if exName == "" {
		panic(errors.New("exName required"))
	}
	if exType == "" {
		panic(errors.New("exType required"))
	}
	if len(routingKeys) == 0 {
		panic(errors.New("routingKeys required"))
	}
	if !consume {
		if contentType == "" {
			panic(errors.New("contentType required"))
		}
	}

	return &QC{
		url:         url,
		exName:      exName,
		exType:      exType,
		routingKeys: routingKeys,
		consume:     consume,
		contentType: contentType,
		durable:     durable,
	}
}

func (qc *QC) Connect() error {

	conn, err := amqp.Dial(qc.url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	qc.conn = conn
	qc.ch = ch

	return nil
}

func (qc *QC) declareExchange() error {

	err := qc.ch.ExchangeDeclare(
		qc.exName,
		qc.exType,
		qc.durable, // durable
		true,       // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func (qc *QC) declareQueue() (string, error) {
	q, err := qc.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return "", err
	}

	return q.Name, nil
}

func (qc *QC) Run() error {

	defer qc.conn.Close()
	defer qc.ch.Close()

	err := qc.declareExchange()
	if err != nil {
		return err
	}

	if qc.consume {
		qName, err := qc.declareQueue()
		if err != nil {
			return err
		}

		err = qc.bindQueue(qName)
		if err != nil {
			return err
		}

		err = qc.doConsume(qName)
		if err != nil {
			return err
		}

	} else {
		err = qc.doPublish()
		if err != nil {
			return err
		}
	}

	return nil
}

func (qc *QC) doPublish() error {

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	for _, routingKey := range qc.routingKeys {
		err = qc.ch.Publish(
			qc.exName, // exchange
			routingKey,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: qc.contentType,
				Body:        bytes,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (qc *QC) bindQueue(qName string) error {

	var err error

	for _, routingKey := range qc.routingKeys {
		err = qc.ch.QueueBind(
			qName,
			routingKey,
			qc.exName,
			false, // no-wait
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (qc *QC) doConsume(qName string) error {

	msgs, err := qc.ch.Consume(
		qName, // queue
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		os.Stdout.Write(d.Body)
	}

	return nil
}
