package handlers

import (
	"encoding/json"
	"net/http"
)

func (dbCfg *DBConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request, JWTSecret string) {
	type SuccessRes struct {
		Body string `json:"cleaned_body"`
	}

	isAuthenticated, userId := dbCfg.isAuthenticated(w, r, JWTSecret)

	if !isAuthenticated {
		return
	}

	decoder := json.NewDecoder(r.Body)
	body := ResBody{}
	err := decoder.Decode(&body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	nonProfaneMsg := removeProfanity(body.Body)

	chirp, err := dbCfg.DB.CreateChirp(nonProfaneMsg, userId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, 201, chirp)
}

func (dbCfg *DBConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := dbCfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, 200, chirps)
}

func (dbCfg *DBConfig) GetSingleChirpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	chirp, err := dbCfg.DB.GetSingleChirp(id)

	if err != nil {
		if err.Error() == "chirp not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, 200, chirp)
}
