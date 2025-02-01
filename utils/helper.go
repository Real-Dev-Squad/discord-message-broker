package utils

import (
	"bytes"
	"errors"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/Real-Dev-Squad/discord-message-broker/config"
	"github.com/sirupsen/logrus"
)

var ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		logrus.Errorf("Attempt %d: Operation failed: %s", i+1, err)
		if i < maxRetries-1 {
			backoffDuration := time.Duration(math.Pow(2, float64(i))) * time.Second

			if backoffDuration > 30*time.Second {
				backoffDuration = 30 * time.Second
			}
			time.Sleep(backoffDuration)

		}
	}
	return err
}

var MakeAPIRequest = func(method string, endpoint string, body *[]byte) (*[]byte, error) {
	client := &http.Client{
		Timeout: config.AppConfig.API_TIMEOUT,
	}

	newBody := body
	if newBody == nil {
		newBody = &[]byte{}
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(*newBody))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get a successful response from the API")
	}

	return &bodyBytes, nil
}

var SendDataToDiscordService = func(body []byte) error {
	endpoint := config.AppConfig.DISCORD_SERVICE_URL + config.AppConfig.DISCORD_SERVICE_ENDPOINT
	return ExponentialBackoffRetry(5, func() error {
		responseBody, err := MakeAPIRequest("POST", endpoint, &body)
		if err != nil {
			return err
		}
		logrus.Infof("Received response from API: %s", string(*responseBody))
		return nil
	})

}
