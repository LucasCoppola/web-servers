package database

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
