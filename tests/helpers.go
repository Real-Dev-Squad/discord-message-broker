package helpers

import (
	"os"
)

func init() {
	os.Setenv("QUEUE_NAME", "DISCORD_QUEUE")
	os.Setenv("QUEUE_URL", "local://123")
}
