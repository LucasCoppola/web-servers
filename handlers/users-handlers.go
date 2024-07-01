package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (dbCfg *DBConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var usrBody UserBody
	err := decoder.Decode(&usrBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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

func (dbCfg *DBConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var usrBody UserBody
	err := decoder.Decode(&usrBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userFound, exists, err := dbCfg.DB.FindUserByEmail(usrBody.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		respondWithError(w, http.StatusBadRequest, "User doesn't exists")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(usrBody.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User doesn't exists")
		return
	}

	userResponse := UserResponse{Id: userFound.Id, Email: userFound.Email}

	respondWithJSON(w, 200, userResponse)
}
