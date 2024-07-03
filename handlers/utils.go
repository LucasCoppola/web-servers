package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func CreateJWT(exp_in_secs int, userId int, JWTSecret string) (string, error) {
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(exp_in_secs) * time.Second)),
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
