package log

import (
	"fmt"
	"os"
	"time"

	"github.com/rkriad585/neovector/internal/config"
)

func Write(entry string) error {
	if err := config.EnsureConfigDir(); err != nil {
		return err
	}
	path := config.LogPath()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	_, err = fmt.Fprintf(file, "[%s] %s\n", timestamp, entry)
	return err
}
