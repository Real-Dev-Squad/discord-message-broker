package utils

import (
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

// TODO: Remove this, and use the one from github.com/Real-Dev-Squad/discord-service/utils/helper
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
