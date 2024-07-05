package database

import (
	"errors"
	"sort"
	"strconv"
)

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, userId int) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStructure.Chirps) + 1
	chirp := Chirp{Id: newId, Body: body, AuthorId: userId}

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
func (db *DB) GetChirps(authorId string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return []Chirp{}, err
	}

	var chirpsList []Chirp

	if authorId == "" {
		for _, chirp := range dbStructure.Chirps {
			chirpsList = append(chirpsList, chirp)
		}
	} else {
		authorIdInt, err := strconv.Atoi(authorId)
		if err != nil {
			return nil, err
		}
		for _, chirp := range dbStructure.Chirps {
			if chirp.AuthorId == authorIdInt {
				chirpsList = append(chirpsList, chirp)
			}
		}
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

func (db *DB) DeleteChirp(id string, userId int) (int, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return 500, err
	}

	parsedId, err := strconv.Atoi(id)

	if err != nil {
		return 500, err
	}

	chirp, ok := dbStructure.Chirps[parsedId]

	if !ok {
		return 404, errors.New("chirp not found")
	}

	if chirp.AuthorId != userId {
		return 403, errors.New("You're not the author of the chirp")
	}

	delete(dbStructure.Chirps, parsedId)

	return 204, nil
}
