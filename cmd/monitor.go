/* github.com/rhinocodelab/fic/cmd/monitor.go */

package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rhinocodelab/fic/internal/config"
	"github.com/rhinocodelab/fic/internal/hasher"
	"github.com/rhinocodelab/fic/internal/logger"
	"github.com/spf13/cobra"
)

const (
	tempDBPath = "/tmp/temp_db.json"
)

type FileHashEntry struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// FileHash structure for tracking file paths and hashes
type FileHash struct {
	FilePath map[string]FileHashEntry `json:"files"`
}

// createTemporaryDB generates a JSON file with current file hashes
func createTemporaryDB(scanPaths []string) error {
	tempDB := FileHash{FilePath: make(map[string]FileHashEntry)}

	for _, path := range scanPaths {
		err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				logger.Logger.Printf("Error accessing path %s: %v", filePath, err)
				return nil
			}

			if !d.IsDir() {
				hash, err := hasher.CalculateHash(filePath)
				if err != nil {
					logger.Logger.Printf("Failed to hash file %s: %v", filePath, err)
					return nil
				}
				tempDB.FilePath[filePath] = FileHashEntry{Path: filePath, Hash: hash}
			}
			return nil
		})

		if err != nil {
			logger.Logger.Printf("Error scanning path %s: %v", path, err)
			continue
		}
	}

	data, err := json.MarshalIndent(tempDB, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tempDB: %v", err)
	}

	err = os.WriteFile(tempDBPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tempDB file: %v", err)
	}

	logger.Logger.Println("Temporary database created successfully at", tempDBPath)
	return nil
}

// Database structure
type Database struct {
	Files      map[string]FileHashEntry `json:"files"`
	TotalFiles int                      `json:"total_files"`
	DBPath     string                   `json:"db_path"`
}

// compareDatabases compares the temporary DB with the actual DB
func compareDatabases(actualDBPath string) {
	logger.Logger.Println("Comparing temporary DB with actual DB...")

	actualDBFile, err := os.ReadFile(actualDBPath)
	if err != nil {
		logger.Logger.Println("Error reading actual database:", err)
		return
	}

	tempDBFile, err := os.ReadFile(tempDBPath)
	if err != nil {
		logger.Logger.Println("Error reading temporary database:", err)
		return
	}

	var actualDB, tempDB Database

	if err := json.Unmarshal(actualDBFile, &actualDB); err != nil {
		logger.Logger.Println("Error parsing actual database:", err)
		return
	}

	if err := json.Unmarshal(tempDBFile, &tempDB); err != nil {
		logger.Logger.Println("Error parsing temporary database:", err)
		return
	}

	// Track changes
	for path, tempEntry := range tempDB.Files {
		actualEntry, exists := actualDB.Files[path]

		if exists {
			// File exists — check hash
			if actualEntry.Hash == tempEntry.Hash {
				logger.Logger.Printf("Hash matched for: %s", path)
			} else {
				logger.Logger.Printf("Hash mismatch for: %s", path)
			}
		} else {
			// New file detected
			logger.Logger.Printf("New file detected: %s", path)
		}
	}

	// Check for deleted files
	for path := range actualDB.Files {
		if _, exists := tempDB.Files[path]; !exists {
			logger.Logger.Printf("File deleted: %s", path)
		}
	}
}

func runDaemon(cfg *config.Config) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Logger.Println("Starting fic monitor daemon...")

	for {
		startTime := time.Now()

		if err := createTemporaryDB(cfg.ScanPaths); err != nil {
			logger.Logger.Println("Error creating temporary database:", err)
		} else {
			compareDatabases(cfg.DatabasePath)
		}

		elapsed := time.Since(startTime)
		logger.Logger.Printf("Monitor cycle completed in %v", elapsed)

		// Ensure 1-minute interval between cycles
		time.Sleep(time.Until(startTime.Add(1 * time.Minute)))
	}

	// Wait for termination signal
	sig := <-signalChan
	logger.Logger.Printf("Received signal: %v. Exiting...", sig)
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor file changes and update HASH values",
	Long:  `Monitor the files under scan path`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Monitoring files for changes...")

		// Load configuration
		cfg, err := config.LoadConfig("config/config.json")
		if err != nil {
			logger.Logger.Println("Error loading config:", err)
			return
		}
		// Run as daemon
		runDaemon(cfg)
		// Keep main goroutine running
		select {}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
