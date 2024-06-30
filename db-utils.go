package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
func NewDB(path string) (*DB, error) {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	err := db.ensureDB()

	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStructure.Chirps) + 1
	chirp := Chirp{Id: newId, Body: body}

	if dbStructure.Chirps == nil {
		dbStructure.Chirps = make(map[int]Chirp)
	}
	dbStructure.Chirps[newId] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.loadDB()

	if err != nil {
		return []Chirp{}, err
	}

	var chirpsList []Chirp

	for _, chirp := range data.Chirps {
		chirpsList = append(chirpsList, Chirp{Id: chirp.Id, Body: chirp.Body})
	}

	return chirpsList, nil
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
			return DBStructure{Chirps: make(map[int]Chirp)}, nil
		}
		return DBStructure{}, err
	}

	if len(data) == 0 {
		// If the file is empty, return an empty DBStructure
		return DBStructure{Chirps: make(map[int]Chirp)}, nil
	}

	var dbStructure DBStructure
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	if dbStructure.Chirps == nil {
		dbStructure.Chirps = make(map[int]Chirp)
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
