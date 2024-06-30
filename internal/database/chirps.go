package database

import (
	"errors"
	"sort"
	"strconv"
)

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

	sort.Slice(chirpsList, func(i, j int) bool {
		return chirpsList[i].Id < chirpsList[j].Id
	})

	return chirpsList, nil
}

// GetSingleChirp returns a single chirp from the database
func (db *DB) GetSingleChirp(id string) (Chirp, error) {
	data, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	parsedId, err := strconv.Atoi(id)

	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := data.Chirps[parsedId]

	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}

	return chirp, nil
}
