package handlers

import (
	"encoding/json"
	"net/http"
)

func (dbCfg *DBConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type userBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	var usrBody userBody
	err := decoder.Decode(&usrBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	user, err := dbCfg.DB.CreateUser(usrBody.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, 201, user)
}

func (dbCfg *DBConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {}
