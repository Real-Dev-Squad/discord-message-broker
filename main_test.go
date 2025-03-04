package main

import (
	"errors"
	"testing"

	"github.com/Real-Dev-Squad/discord-message-broker/config"
	"github.com/Real-Dev-Squad/discord-message-broker/model"
	_ "github.com/Real-Dev-Squad/discord-message-broker/tests"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

type mockQueue struct {
	dialError    error
	channelError error
	queueError   error
}

func (m *mockQueue) Consumer() {
}

func (m *mockQueue) Dial() error {
	return m.dialError
}

func (m *mockQueue) CreateChannel() error {
	return m.channelError
}
func (m *mockQueue) DeclareQueue() error {
	return m.queueError
}

func TestInitQueueConnection(t *testing.T) {
	config.AppConfig.MAX_RETRIES = 1
	t.Run("should not panic when Dial() returns error", func(t *testing.T) {
		mockQueue := &mockQueue{dialError: errors.New("connection failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not panic when Dial is unsuccessful")

	})

	t.Run("should not panic when CreateChannel() returns error", func(t *testing.T) {
		mockQueue := &mockQueue{channelError: errors.New("channel failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not panic when CreateChannel is unsuccessful")

	})

	t.Run("should not panic when DeclareQueue() returns error", func(t *testing.T) {
		mockQueue := &mockQueue{queueError: errors.New("queue failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not when DeclareQueue is unsuccessful")

	})

}

func TestSessionWrapper(t *testing.T) {
	sessionWrapper := &model.Queue{}

	t.Run("SessionWrapper should always implement dial() method", func(t *testing.T) {
		err := sessionWrapper.Dial()
		assert.Error(t, err)
	})

	t.Run("SessionWrapper should always implement createChannel() method", func(t *testing.T) {
		sessionWrapper.Connection = &amqp.Connection{}
		assert.Panics(t, func() {
			sessionWrapper.CreateChannel()
		})

	})

	t.Run("SessionWrapper should always implement declareQueue() method", func(t *testing.T) {
		sessionWrapper.Channel = &amqp.Channel{}
		assert.Panics(t, func() {
			sessionWrapper.DeclareQueue()
		})
	})

}
