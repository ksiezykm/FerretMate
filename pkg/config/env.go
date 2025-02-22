package config

import (
	"bufio"
	"fmt"
	"os"
)

// readDBURI reads the MongoDB URI from the .env file
func ReadDBURI(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", fmt.Errorf("empty or invalid .env file")
}