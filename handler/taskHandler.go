package handler

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func TaskHandler(
	msgs <-chan amqp.Delivery,
) {

	for d := range msgs {
		logrus.Printf("Received a message: %s", d.Body)
	}

}
