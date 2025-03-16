package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DatabasePath string   `json:"database_path"`
	LogFilePath  string   `json:"fic_log_path"`
	ScanPaths    []string `json:"scan_paths"`
}

// LoadConfig reads and parses the config.json file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
