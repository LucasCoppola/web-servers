package database

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chirp, err := db.CreateUser("jhondoe@gmail.com")
	if err != nil {
		t.Fatalf("Failed to create chirp: %v", err)
	}

	if chirp.Id != 1 {
		t.Errorf("Expected chirp ID to be 1, got %d", chirp.Id)
	}

	if chirp.Email != "jhondoe@gmail.com" {
		t.Errorf("Expected chirp body to be 'jhondoe@gmail.com', got '%s'", chirp.Email)
	}
}
