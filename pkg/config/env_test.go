package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReadDBURI(t *testing.T) {
	// Test case 1: File with a URI
	t.Run("file with a URI", func(t *testing.T) {
		tempDir := t.TempDir()
		expectedURI := "mongodb://localhost:27017/testdb"
		envFilePath := filepath.Join(tempDir, ".env")

		if err := os.WriteFile(envFilePath, []byte(expectedURI), 0600); err != nil {
			t.Fatalf("Failed to write test .env file: %v", err)
		}

		uri, err := ReadDBURI(envFilePath)
		if err != nil {
			t.Errorf("ReadDBURI() error = %v, wantErr %v", err, false)
		}
		if uri != expectedURI {
			t.Errorf("ReadDBURI() gotURI = %q, want %q", uri, expectedURI)
		}
	})

	// Test case 2: Empty file
	t.Run("empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		envFilePath := filepath.Join(tempDir, ".env")

		if err := os.WriteFile(envFilePath, []byte(""), 0600); err != nil {
			t.Fatalf("Failed to write empty test .env file: %v", err)
		}

		_, err := ReadDBURI(envFilePath)
		expectedErr := fmt.Errorf("empty or invalid .env file")
		if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("ReadDBURI() error = %v, want %v", err, expectedErr)
		}
	})

	// Test case 3: File not found
	t.Run("file not found", func(t *testing.T) {
		_, err := ReadDBURI("non_existent_env_file.env")
		if !os.IsNotExist(err) {
			t.Errorf("ReadDBURI() error = %v, want os.ErrNotExist", err)
		}
	})

	// Test case 4: File with multiple lines
	t.Run("file with multiple lines", func(t *testing.T) {
		tempDir := t.TempDir()
		firstLineURI := "mongodb://user:pass@host:port/dbname"
		multiLineContent := firstLineURI + "\nSECOND_LINE_VAR=some_value\nTHIRD_LINE_VAR=another_value"
		envFilePath := filepath.Join(tempDir, ".env")

		if err := os.WriteFile(envFilePath, []byte(multiLineContent), 0600); err != nil {
			t.Fatalf("Failed to write multi-line test .env file: %v", err)
		}

		uri, err := ReadDBURI(envFilePath)
		if err != nil {
			t.Errorf("ReadDBURI() error = %v, wantErr %v", err, false)
		}
		if uri != firstLineURI {
			t.Errorf("ReadDBURI() gotURI = %q, want %q", uri, firstLineURI)
		}
	})
}
