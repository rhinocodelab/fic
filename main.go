/* github.com/rhinocodelab/fic/main.go */

package main

import (
	"log"
	"os"

	"github.com/rhinocodelab/fic/cmd"
	"github.com/rhinocodelab/fic/internal/config"
	"github.com/rhinocodelab/fic/internal/database"
	"github.com/rhinocodelab/fic/internal/logger"
)

func main() {
	// Load config.json
	cfg, err := config.LoadConfig("config/config.json")
	handleError(err, "Failed to load config", nil)

	// Initialize logger
	err = logger.InitLogger(cfg.LogFilePath)
	handleError(err, "Failed to initialize logger", nil)

	// Create Base database
	db := database.NewDatabase(cfg.DatabasePath)
	err = db.CreateBaseDB()
	handleError(err, "Failed to create base database", logger.Logger)

	cmd.Execute()
}

func handleError(err error, message string, logger *log.Logger) {
	if err != nil {
		if logger != nil {
			logger.Fatalf("%s: %v", message, err)
		} else {
			log.Fatalf("%s: %v", message, err)
		}
		os.Exit(1)
	}
}
