package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	QUEUE_NAME               string
	MAX_RETRIES              int
	QUEUE_URL                string
	DISCORD_SERVICE_URL      string
	API_TIMEOUT              time.Duration
	DISCORD_SERVICE_ENDPOINT string
}

var AppConfig Config

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("Loaded .env file successfully")
	}

	AppConfig = Config{
		QUEUE_NAME:               loadEnv("QUEUE_NAME"),
		QUEUE_URL:                loadEnv("QUEUE_URL"),
		DISCORD_SERVICE_URL:      loadEnv("DISCORD_SERVICE_URL"),
		DISCORD_SERVICE_ENDPOINT: "/queue",
		API_TIMEOUT:              5 * time.Minute,
		MAX_RETRIES:              5,
	}
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}
