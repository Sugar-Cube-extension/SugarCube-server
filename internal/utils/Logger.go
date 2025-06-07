package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	logDir := "/var/log/sugarcube-backend"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102-150405")
	filePath := filepath.Join(logDir, fmt.Sprintf("sugarcube-backend-%s.log", timestamp))

	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}

	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()

	log.Logger = log.Output(zerolog.MultiLevelWriter(os.Stdout, logFile))

	log.Info().Msg("Logger initialized successfully.")
}
