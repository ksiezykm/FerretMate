package model

import (
	"encoding/json"
	"os"
)

type Connection struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func LoadConnections() ([]Connection, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var connections []Connection
	if err := json.Unmarshal(data, &connections); err != nil {
		return nil, err
	}

	return connections, nil
}
