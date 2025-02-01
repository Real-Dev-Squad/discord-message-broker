package model

import (
	"github.com/Real-Dev-Squad/discord-message-broker/config"
	"github.com/Real-Dev-Squad/discord-message-broker/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	Connection *amqp.Connection
	Queue      amqp.Queue
	Name       string
	Channel    *amqp.Channel
}

func (q *Queue) Dial() error {
	var err error
	q.Connection, err = amqp.Dial(config.AppConfig.QUEUE_URL)
	return err
}

func (q *Queue) CreateChannel() error {
	var err error
	q.Channel, err = q.Connection.Channel()
	return err
}

func (q *Queue) DeclareQueue() error {
	var err error
	q.Queue, err = q.Channel.QueueDeclare(
		config.AppConfig.QUEUE_NAME,     // name
		true,                            // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		amqp.Table{"x-max-priority": 2}, // arguments
	)
	q.Name = config.AppConfig.QUEUE_NAME // Ensure the queue name is set
	return err
}

func (q *Queue) Consumer() {
	msgs, err := q.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Errorf("%s Failed to register a consumer", err)
		return
	}

	forever := make(chan bool)
	go func() {
		logrus.Info("Consumer connected")
		for d := range msgs {
			logrus.Printf("Received a message: %s", d.Body)
			utils.SendDataToDiscordService(d.Body)
			d.Ack(false)
		}
	}()

	<-forever
	logrus.Info("Consumer stopped")
}

type SessionInterface interface {
	Dial() error
	CreateChannel() error
	DeclareQueue() error
	Consumer()
}
