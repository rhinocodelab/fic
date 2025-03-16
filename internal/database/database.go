/* github.com/rhinocodelab/fic/internal/database/database.go */

package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syscall"

	"github.com/rhinocodelab/fic/internal/logger"
)

type Entry struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Database struct {
	Files      map[string]Entry `json:"files"`
	TotalFiles int              `json:"total_files"`
	DBPath     string           `json:"db_path"`
}

// Initialize a new database
func NewDatabase(path string) *Database {
	return &Database{
		Files:      make(map[string]Entry),
		TotalFiles: 0,
		DBPath:     path,
	}
}

// Create Base database
func (db *Database) CreateBaseDB() error {
	// Check if the database file already exists
	if _, err := os.Stat(db.DBPath); err == nil {
		return nil
	}
	// Create a new database
	if err := os.MkdirAll(filepath.Dir(db.DBPath), 0755); err != nil {
		logger.Logger.Printf("Error creating directory for database: %v", err)
		return err
	}
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		logger.Logger.Printf("Error serializing database: %v", err)
		return err
	}
	if err := os.WriteFile(db.DBPath, data, 0644); err != nil {
		logger.Logger.Printf("Error writing database to file: %v", err)
		return err
	}
	logger.Logger.Printf("Database created successfully at %s", db.DBPath)
	return nil
}

// Load loads the database from the file specified by DBPath. If the file does not exist,
// it initializes an empty database.
func (db *Database) Load() error {
	lockFile, err := db.lockFile()
	if err != nil {
		return err
	}
	defer db.unlockFile(lockFile)
	if _, err := os.Stat(db.DBPath); os.IsNotExist(err) {
		db.Files = make(map[string]Entry)
		db.TotalFiles = 0
		return nil
	}
	data, err := os.ReadFile(db.DBPath)
	if err != nil {
		return err
	}
	// Unmarshal the JSON data into the database struct
	if err := json.Unmarshal(data, db); err != nil {
		return err
	}
	// Ensure TotalFiles is consistent with the loaded data
	db.TotalFiles = len(db.Files)
	return nil
}

// UpdateDB updates the database with a file path and its hash value, adding the entry if it does not exist.
func (db *Database) UpdateDB(filePath, hash string) error {
	// Check if the file already exists in the database
	if entry, exists := db.Files[filePath]; exists && entry.Hash == hash {
		return nil
	}
	// Add the file to the database
	db.Files[filePath] = Entry{Path: filePath, Hash: hash}
	db.TotalFiles = len(db.Files)
	if err := db.Save(); err != nil {
		return err
	}
	// Save the updated database
	return nil
}

// Save persists the current state of the database to the file specified by DBPath.
func (db *Database) Save() error {
	lockFile, err := db.lockFile()
	if err != nil {
		return err
	}
	defer db.unlockFile(lockFile)
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(db.DBPath, data, 0644); err != nil {
		return err
	}
	return nil
}

// lockFile acquires an exclusive lock on the database file to prevent concurrent access.
func (db *Database) lockFile() (*os.File, error) {
	lockFilePath := db.DBPath + ".lock"
	// Create a lock file
	lockFile, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	// Check if the lock file already exists
	if err != nil {
		return nil, err
	}
	// Lock the file
	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		lockFile.Close()
		return nil, err
	}
	// Return the lock file
	return lockFile, nil
}

// unlockFile releases the lock on the database file.
func (db *Database) unlockFile(lockFile *os.File) {
	// Unlock the file
	syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
	// Close the lock file
	lockFile.Close()
	// Remove the lock file
	os.Remove(lockFile.Name())
}
