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

func TestCreateChirp(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chirp, err := db.CreateChirp("Test chirp")
	if err != nil {
		t.Fatalf("Failed to create chirp: %v", err)
	}

	if chirp.Id != 1 {
		t.Errorf("Expected chirp ID to be 1, got %d", chirp.Id)
	}

	if chirp.Body != "Test chirp" {
		t.Errorf("Expected chirp body to be 'Test chirp', got '%s'", chirp.Body)
	}
}

func TestGetChirps(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.CreateChirp("First chirp")
	if err != nil {
		t.Fatalf("Failed to create chirp: %v", err)
	}
	_, err = db.CreateChirp("Second chirp")
	if err != nil {
		t.Fatalf("Failed to create chirp: %v", err)
	}

	chirps, err := db.GetChirps()
	if err != nil {
		t.Fatalf("Failed to get chirps: %v", err)
	}

	if len(chirps) != 2 {
		t.Errorf("Expected 2 chirps, got %d", len(chirps))
	}

	if chirps[0].Body != "First chirp" {
		t.Errorf("Expected first chirp body to be 'First chirp', got '%s'", chirps[0].Body)
	}

	if chirps[1].Body != "Second chirp" {
		t.Errorf("Expected second chirp body to be 'Second chirp', got '%s'", chirps[1].Body)
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
