package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserBody struct {
	Email              string `json:"email"`
	Password           string `json:"password"`
	Expires_in_seconds int    `json:"expires_in_seconds"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

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

	bearerToken := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(bearerToken, "Bearer ")

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	userId, err := strconv.Atoi(claims.Subject)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid user ID format")
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

func (dbCfg *DBConfig) LoginHandler(w http.ResponseWriter, r *http.Request, JWTSecret string) {
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

	var exp_in_secs int
	if usrBody.Expires_in_seconds == 0 {
		exp_in_secs = 86400 // 24 hours in seconds
	} else {
		exp_in_secs = usrBody.Expires_in_seconds
	}

	token, err := CreateJWT(exp_in_secs, userResponse.Id, JWTSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = dbCfg.DB.StoreToken(userFound.Id, token)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse.Token = token

	respondWithJSON(w, 200, userResponse)
}
