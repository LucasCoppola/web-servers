package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	HOUR_IN_SECS = 3600
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorRes struct {
		Error string `json:"error"`
	}

	error := ErrorRes{Error: msg}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(error)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func removeProfanity(msg string) string {
	profaneWords := [3]string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(msg, " ")

	for i, word := range words {
		for _, profaneWord := range profaneWords {
			if profaneWord == strings.ToLower(word) {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}

func CreateJWT(userId int, JWTSecret string) (string, error) {
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(HOUR_IN_SECS) * time.Second)),
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(JWTSecret))

	if err != nil {
		return "", err
	}

	return signedString, nil
}

func generateRandomString() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (dbCfg *DBConfig) isAuthenticated(w http.ResponseWriter, r *http.Request, JWTSecret string) (bool, int) {
	bearerToken := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(bearerToken, "Bearer ")

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return false, 0
	}

	if !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return false, 0
	}

	userId, err := strconv.Atoi(claims.Subject)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid user ID format")
		return false, 0
	}

	return true, userId
}
