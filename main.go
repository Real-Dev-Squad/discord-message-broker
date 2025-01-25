package main

import (
	"github.com/Real-Dev-Squad/discord-message-broker/config"
	"github.com/Real-Dev-Squad/discord-message-broker/model"
	"github.com/Real-Dev-Squad/discord-message-broker/utils"
	"github.com/sirupsen/logrus"
)

func InitQueueConnection(openSession model.SessionInterface) {
	var err error
	f := func() error {
		err = openSession.Dial()
		if err != nil {
			return err
		}
		err = openSession.CreateChannel()
		if err != nil {
			return err
		}
		err = openSession.DeclareQueue()
		return err
	}

	err = utils.ExponentialBackoffRetry(config.AppConfig.MAX_RETRIES, f)
	if err != nil {
		logrus.Errorf("Failed to initialize queue after %d attempts: %s", config.AppConfig.MAX_RETRIES, err)
		return
	}
	logrus.Infof("Established a connection to RabbitMQ named %s", config.AppConfig.QUEUE_NAME)
	openSession.Consumer()
}

func main() {
	queueInstance := &model.Queue{}
	InitQueueConnection(queueInstance)
}
