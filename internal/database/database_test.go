package database

import (
	"os"
	"testing"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "test_db_*.json")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	tempFile.Close()

	db, err := NewDB(tempFile.Name())
	if err != nil {
		os.Remove(tempFile.Name())
		t.Fatalf("Could not create database: %v", err)
	}

	return db, func() {
		os.Remove(tempFile.Name())
	}
}

func TestLoadDB(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.CreateChirp("Test chirp")
	if err != nil {
		t.Fatalf("Failed to create chirp: %v", err)
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		t.Fatalf("Failed to load database: %v", err)
	}

	if len(dbStructure.Chirps) != 1 {
		t.Errorf("Expected 1 chirp in loaded database, got %d", len(dbStructure.Chirps))
	}

	if dbStructure.Chirps[1].Body != "Test chirp" {
		t.Errorf("Expected chirp body to be 'Test chirp', got '%s'", dbStructure.Chirps[1].Body)
	}
}
