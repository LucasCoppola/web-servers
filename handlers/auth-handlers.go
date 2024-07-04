package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

	token, err := CreateJWT(userResponse.Id, JWTSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse.Token = token

	refreshToken, err := generateRandomString()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour).Unix()

	err = dbCfg.DB.StoreRefreshToken(refreshToken, userFound.Id, expiresAt)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse.RefreshToken = refreshToken

	respondWithJSON(w, 200, userResponse)
}

func (dbCfg *DBConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request, JWTSecret string) {
	type TokenResponse struct {
		Token string `json:"token"`
	}

	bearerToken := r.Header.Get("Authorization")
	refreshToken := strings.TrimPrefix(bearerToken, "Bearer ")

	exists, expiresAt, userId, err := dbCfg.DB.FindRefreshToken(refreshToken)

	if err != nil || !exists {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if time.Now().Unix() > expiresAt {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired")
		return
	}

	newAccessToken, err := CreateJWT(userId, JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate new access token")
		return
	}

	token := TokenResponse{Token: newAccessToken}

	respondWithJSON(w, 200, token)
}

func (dbCfg *DBConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request, JWTSecret string) {
	bearerToken := r.Header.Get("Authorization")
	refreshToken := strings.TrimPrefix(bearerToken, "Bearer ")

	exists, _, userId, err := dbCfg.DB.FindRefreshToken(refreshToken)

	if err != nil || !exists {
		respondWithError(w, http.StatusBadRequest, "Invalid refresh token")
		return
	}

	err = dbCfg.DB.RevokeRefreshToken(userId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token")
		return
	}

	w.WriteHeader(204)
}
