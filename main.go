package main

import (
	"sync"

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

func (q *Queue) dial() error {
	var err error
	q.Connection, err = amqp.Dial(config.AppConfig.QUEUE_URL)
	return err
}

func (q *Queue) createChannel() error {
	var err error
	q.Channel, err = q.Connection.Channel()
	return err
}

func (q *Queue) declareQueue() error {
	var err error
	q.Queue, err = q.Channel.QueueDeclare(
		config.AppConfig.QUEUE_NAME,     // name
		true,                            // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		amqp.Table{"x-max-priority": 2}, // arguments
	)
	return err
}

func (q *Queue) consumer() {
	msgs, err := q.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Errorf("%s Failed to register a consumer", err)
	}

	//TODO: Implement API with authentication (tracking issue: https://github.com/Real-Dev-Squad/discord-service/issues/28)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			logrus.Printf("Received a message: %s", d.Body)
		}
	}()

	<-forever
}

var (
	queueInstance *Queue
	once          sync.Once
)

type sessionInterface interface {
	dial() error
	createChannel() error
	declareQueue() error
	consumer()
}

func InitQueueConnection(openSession sessionInterface) {
	var err error
	f := func() error {
		err = openSession.dial()
		if err != nil {
			return err
		}
		err = openSession.createChannel()
		if err != nil {
			return err
		}
		err = openSession.declareQueue()
		return err
	}

	err = utils.ExponentialBackoffRetry(config.AppConfig.MAX_RETRIES, f)
	if err != nil {
		logrus.Errorf("Failed to initialize queue after %d attempts: %s", config.AppConfig.MAX_RETRIES, err)
		return
	}
	openSession.consumer()
	logrus.Infof("Established a connection to RabbitMQ named %s", config.AppConfig.QUEUE_NAME)

}

func main() {
	queueInstance := &Queue{}
	InitQueueConnection(queueInstance)
}
