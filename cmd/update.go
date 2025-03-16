/* github.com/rhinocodelab/fic/cmd/update.go */

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rhinocodelab/fic/internal/config"
	"github.com/rhinocodelab/fic/internal/database"
	"github.com/rhinocodelab/fic/internal/hasher"
	"github.com/rhinocodelab/fic/internal/logger"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the database with new SHA256 Hash",
	Long:  "Calculate the SHA256 Hash of the scan path files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating SHA256 Hash")

		// Load configuration
		cfg, err := config.LoadConfig("config/config.json")
		if err != nil {
			logger.Logger.Println("Error loading config:", err)
			return
		}

		// Load Database
		db := database.NewDatabase(cfg.DatabasePath)
		err = db.CreateBaseDB()
		if err != nil {
			logger.Logger.Println("Error creating database:", err)
			return
		}
		// Verify each scan path
		for _, scanPath := range cfg.ScanPaths {
			// Check if the scan path exists
			if _, err := os.Stat(scanPath); os.IsNotExist(err) {
				logger.Logger.Printf("Scan path %s does not exist\n", scanPath)
				continue
			}
			files, err := os.ReadDir(scanPath)
			if err != nil {
				logger.Logger.Printf("Error reading scan path %s: %v\n", scanPath, err)
				continue
			}

			for _, file := range files {
				fullPath := filepath.Join(scanPath, file.Name())
				hash, err := hasher.CalculateHash(fullPath)
				if err != nil {
					logger.Logger.Printf("Error calculating hash for %s: %v\n", fullPath, err)
					continue
				}
				// Update the database with the new hash
				err = db.UpdateDB(fullPath, hash)
				if err != nil {
					logger.Logger.Printf("Error updating hash for %s: %v\n", fullPath, err)
					continue
				}

			}

		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
