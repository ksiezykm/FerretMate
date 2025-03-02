package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Config struct {
	Databases []DatabaseConfig `json:"databases"`
}

func loadConfig(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	return config, err
}

func readConfig() (Config, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "ferretmate")
	configFile := filepath.Join(configDir, "config.json")

	config, err := loadConfig(configFile)
	if err != nil {
		return config, fmt.Errorf("Error readConfig: %w", err)
	}
	return config, nil
}
