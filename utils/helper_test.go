package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Real-Dev-Squad/discord-message-broker/config"
	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoffRetry_Success(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := ExponentialBackoffRetry(5, operation)
	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestExponentialBackoffRetry_Failure(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("permanent error")
	}

	err := ExponentialBackoffRetry(3, operation)
	assert.Error(t, err)
	assert.Equal(t, 3, attempts)
}

func TestExponentialBackoffRetry_NoRetries(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("error")
	}

	err := ExponentialBackoffRetry(0, operation)
	assert.Nil(t, err)
	assert.Equal(t, 0, attempts)
}

func TestExponentialBackoffRetry_ImmediateSuccess(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return nil
	}
	err := ExponentialBackoffRetry(5, operation)
	assert.NoError(t, err)
	assert.Equal(t, 1, attempts)
}

func TestMakeAPIRequest_Success(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer mockServer.Close()

	// Test the MakeAPIRequest function
	method := "GET"
	endpoint := mockServer.URL
	body := []byte{}

	responseBody, err := MakeAPIRequest(method, endpoint, &body)
	assert.NoError(t, err)
	assert.NotNil(t, responseBody)
	assert.Equal(t, "Hello, World!", string(*responseBody))
}

func TestMakeAPIRequest_Failure(t *testing.T) {
	// Create a mock server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Test the MakeAPIRequest function
	method := "GET"
	endpoint := mockServer.URL
	body := []byte{}

	responseBody, err := MakeAPIRequest(method, endpoint, &body)
	assert.Error(t, err)
	assert.Nil(t, responseBody)
	assert.Equal(t, "failed to get a successful response from the API", err.Error())
}

func TestMakeAPIRequest_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	method := "GET"
	endpoint := mockServer.URL
	body := []byte{}

	responseBody, err := MakeAPIRequest(method, endpoint, &body)
	assert.NoError(t, err)
	assert.NotNil(t, responseBody)
}

func TestSendDataToDiscordService(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))
	defer mockServer.Close()

	config.AppConfig.DISCORD_SERVICE_URL = mockServer.URL
	config.AppConfig.DISCORD_SERVICE_ENDPOINT = "/test-endpoint"

	body := []byte(`{"message": "Hello, Discord!"}`)

	err := SendDataToDiscordService(body)
	assert.NoError(t, err)
}
