package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func loadConfigMap(file string) (map[string]DatabaseConfig, error) {
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var configMap map[string]DatabaseConfig
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&configMap)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

func ReadConfig() (map[string]DatabaseConfig, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "ferretmate")
	configFile := filepath.Join(configDir, "config.json")

	config, err := loadConfigMap(configFile)
	if err != nil {
		return config, fmt.Errorf("Error readConfig: %w", err)
	}
	return config, nil
}
