/* github.com/rhinocodelab/fic/internal/logger/logger.go */

package logger

import (
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

// InitLogger initializes the logger
func InitLogger(logFilePath string) error {
	// Check if log file path already exists
	if _, err := os.Stat(logFilePath); err == nil {
		// File exists, open it
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Warning: Failed to open log file, using default logger. Error: %v", err)
			Logger = log.Default()
			return err
		}
		Logger = log.New(file, "FIC: ", log.Ldate|log.Ltime)
		return nil
	}
	// Ensure the log file's parent directory exists
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Warning: Failed to create log directory: %v. Using default logger.", err)
		Logger = log.Default()
		return err
	}

	// Create log file if it doesn't exist
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Warning: Failed to create log file, using default logger. Error: %v", err)
		Logger = log.Default()
		return err
	}

	Logger = log.New(file, "FIC: ", log.Ldate|log.Ltime)
	return nil
}
