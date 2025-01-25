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

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

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

func TestMakeAPIRequest(t *testing.T) {
	t.Run("should return nil when successful", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		}))
		defer mockServer.Close()

		method := "GET"
		endpoint := mockServer.URL
		body := []byte{}

		responseBody, err := MakeAPIRequest(method, endpoint, &body)
		assert.NoError(t, err)
		assert.NotNil(t, responseBody)
		assert.Equal(t, "Hello, World!", string(*responseBody))
	})

	t.Run("should return error when status code is not 200", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		method := "GET"
		endpoint := mockServer.URL
		body := []byte{}

		responseBody, err := MakeAPIRequest(method, endpoint, &body)
		assert.Error(t, err)
		assert.Nil(t, responseBody)
		assert.Equal(t, "failed to get a successful response from the API", err.Error())
	})

	t.Run("should handle timeouts", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		method := "GET"
		endpoint := mockServer.URL
		responseBody, err := MakeAPIRequest(method, endpoint, nil)
		assert.NoError(t, err)
		assert.NotNil(t, responseBody)
	})
	t.Run("should return error if unable to make the request", func(t *testing.T) {
		method := "GET"
		endpoint := "testing"
		_, err := MakeAPIRequest(method, endpoint, nil)
		assert.Error(t, err)
	})

}

func TestSendDataToDiscordService(t *testing.T) {
	t.Run("should return nil when successful", func(t *testing.T) {
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
	})

	t.Run("should return error when unsuccessful", func(t *testing.T) {
		body := []byte(`{"message": "Hello, Discord!"}`)
		err := SendDataToDiscordService(body)
		assert.Error(t, err)
	})

}
