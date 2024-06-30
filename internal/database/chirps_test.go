package database

import (
	"testing"
)

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

func TestGetSingleChirp(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.CreateChirp("First Chirp")
	if err != nil {
		t.Fatalf("Failed to create chirp: %v", err)
	}
	_, err = db.CreateChirp("Second chirp")
	if err != nil {
		t.Fatalf("Failed to create chirp: %v", err)
	}

	chirp, err := db.GetSingleChirp("1")
	if err != nil {
		t.Fatalf("Failed to get chirp: %v", err)
	}

	if chirp.Body != "First Chirp" {
		t.Errorf("Expected chirp body to be 'First chirp', got '%s'", chirp.Body)
	}

	if chirp.Id != 1 {
		t.Errorf("Expected chirp id to be 1, got '%s'", chirp.Body)
	}

	errorChirp, err := db.GetSingleChirp("3")
	if err == nil {
		t.Errorf("Expected error 'chirp not found', got '%v'", errorChirp)
	}
}
