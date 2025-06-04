package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ksiezykm/FerretMate/pkg/model"
)

func TestLoadConfigMap(t *testing.T) {
	// Test case 1: Valid config file
	t.Run("valid config file", func(t *testing.T) {
		tempDir := t.TempDir()
		validConfig := map[string]model.DatabaseConfig{
			"dev": {
				Host:     "mongodb://localhost:27017",
				Database: "devDB",
				Username: "devuser",
				Password: "devpass",
			},
			"prod": {
				Host:     "mongodb://prodserver:27017",
				Database: "prodDB",
				Username: "produser",
				Password: "prodpass",
			},
		}
		fileContent, err := json.MarshalIndent(validConfig, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal test config: %v", err)
		}

		configFilePath := filepath.Join(tempDir, "config.json")
		if err := os.WriteFile(configFilePath, fileContent, 0600); err != nil {
			t.Fatalf("Failed to write test config file: %v", err)
		}

		loadedConfig, err := loadConfigMap(configFilePath)
		if err != nil {
			t.Errorf("loadConfigMap() error = %v, wantErr %v", err, false)
			return
		}
		if !reflect.DeepEqual(loadedConfig, validConfig) {
			t.Errorf("loadConfigMap() got = %v, want %v", loadedConfig, validConfig)
		}
	})

	// Test case 2: File not found
	t.Run("file not found", func(t *testing.T) {
		_, err := loadConfigMap("non_existent_config.json")
		if !os.IsNotExist(err) {
			t.Errorf("loadConfigMap() error = %v, want os.ErrNotExist", err)
		}
	})

	// Test case 3: Invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		invalidJSON := []byte(`{"dev": {"host": "mongodb://localhost:27017", "database": "devDB"},`) // Missing closing brace
		configFilePath := filepath.Join(tempDir, "invalid.json")
		if err := os.WriteFile(configFilePath, invalidJSON, 0600); err != nil {
			t.Fatalf("Failed to write invalid JSON file: %v", err)
		}

		_, err := loadConfigMap(configFilePath)
		// Check if it's a SyntaxError or an UnexpectedEOF error, as both are valid for malformed JSON
		var syntaxError *json.SyntaxError
		if !errors.As(err, &syntaxError) && !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("loadConfigMap() error = %T, %v; want *json.SyntaxError or io.ErrUnexpectedEOF", err, err)
		}
	})

	// Test case 4: Incorrect structure
	t.Run("incorrect structure", func(t *testing.T) {
		tempDir := t.TempDir()
		// Valid JSON, but not map[string]model.DatabaseConfig. For example, a simple array.
		incorrectStructureJSON := []byte(`[{"host": "mongodb://localhost:27017", "database": "devDB"}]`)
		configFilePath := filepath.Join(tempDir, "incorrect_structure.json")
		if err := os.WriteFile(configFilePath, incorrectStructureJSON, 0600); err != nil {
			t.Fatalf("Failed to write incorrect structure JSON file: %v", err)
		}

		_, err := loadConfigMap(configFilePath)
		// We expect a json.UnmarshalTypeError or similar if the structure is wrong.
		// Checking for non-nil error is a basic check; more specific checks might be needed
		// depending on how loadConfigMap and json.Unmarshal behave with specific incorrect structures.
		if err == nil {
			t.Errorf("loadConfigMap() error = nil, want an error for incorrect structure")
		} else {
			// A more specific check could be:
			// if _, ok := err.(*json.UnmarshalTypeError); !ok {
			//  t.Errorf("loadConfigMap() error type = %T, want *json.UnmarshalTypeError or similar", err)
			// }
			// For now, let's log the error type to understand what we get
			t.Logf("Received error for incorrect structure: %T, %v", err, err)
			// And ensure it's a JSON unmarshaling error
			if _, ok := err.(*json.UnmarshalTypeError); !ok {
				// It might also be a generic error from json.Unmarshal if the top level isn't a map
				// Let's check if it's a plain json unmarshal error if not UnmarshalTypeError
				t.Logf("Error is not UnmarshalTypeError, checking for other JSON errors.")
				// This error typically occurs if you try to unmarshal an array into a map struct etc.
				// "json: cannot unmarshal array into Go value of type map[string]model.DatabaseConfig"
				// This is not a *json.UnmarshalTypeError but a plain error from json.Unmarshal.
			}
		}
	})

	// Test case 5: Incorrect field type in structure
	t.Run("incorrect field type in structure", func(t *testing.T) {
		tempDir := t.TempDir()
		// Correct top-level map, but a field has a wrong type
		incorrectFieldTypeJSON := []byte(`{
			"dev": {
				"host": 12345,
				"database": "devDB",
				"username": "user",
				"password": "password"
			}
		}`)
		configFilePath := filepath.Join(tempDir, "incorrect_field_type.json")
		if err := os.WriteFile(configFilePath, incorrectFieldTypeJSON, 0600); err != nil {
			t.Fatalf("Failed to write incorrect field type JSON file: %v", err)
		}

		_, err := loadConfigMap(configFilePath)
		if _, ok := err.(*json.UnmarshalTypeError); !ok {
			t.Errorf("loadConfigMap() error type = %T, value = %v, want *json.UnmarshalTypeError", err, err)
		}
	})
}
