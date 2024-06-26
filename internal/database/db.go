package database

import (
	"encoding/json"
	"os"
	"sync"
)

// NewDB creates a new database connection
func NewDB(path string) (*DB, error) {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	err := db.ensureDB()

	if err != nil {
		return nil, err
	}

	return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()
	_, err := os.Stat(db.path)

	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, create it
			file, err := os.Create(db.path)
			if err != nil {
				return err
			}
			defer file.Close()
		} else {
			return err
		}
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, return an empty DBStructure
			return DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User)}, nil
		}
		return DBStructure{}, err
	}

	if len(data) == 0 {
		// If the file is empty, return an empty DBStructure
		return DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User)}, nil
	}

	var dbStructure DBStructure
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	if dbStructure.Chirps == nil {
		dbStructure.Chirps = make(map[int]Chirp)
	}

	if dbStructure.Users == nil {
		dbStructure.Users = make(map[int]User)
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0644)

	if err != nil {
		return err
	}

	return nil
}
