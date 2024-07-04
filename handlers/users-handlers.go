package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (dbCfg *DBConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var usrBody UserBody
	err := decoder.Decode(&usrBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, exists, err := dbCfg.DB.FindUserByEmail(usrBody.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if exists {
		respondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usrBody.Password), 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := dbCfg.DB.CreateUser(usrBody.Email, hashedPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, 201, user)
}

func (dbCfg *DBConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request, JWTSecret string) {
	decoder := json.NewDecoder(r.Body)
	var usrBody UserBody
	err := decoder.Decode(&usrBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	isAuthenticated, userId := dbCfg.isAuthenticated(w, r, JWTSecret)

	if !isAuthenticated {
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usrBody.Password), 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := dbCfg.DB.UpdateUser(userId, usrBody.Email, hashedPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, 200, user)
}
