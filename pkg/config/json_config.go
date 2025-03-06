package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/ksiezykm/FerretMate/pkg/model"
)


func loadConfigMap(file string) (map[string]model.DatabaseConfig, error) {
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var configMap map[string]model.DatabaseConfig
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&configMap)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

func ReadConfig() (map[string]model.DatabaseConfig, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "ferretmate")
	configFile := filepath.Join(configDir, "config.json")

	config, err := loadConfigMap(configFile)
	if err != nil {
		return config, fmt.Errorf("Error readConfig: %w", err)
	}
	return config, nil
}
