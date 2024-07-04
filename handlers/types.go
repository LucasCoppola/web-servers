package handlers

import (
	"github.com/LucasCoppola/web-server/internal/database"
)

type ApiConfig struct {
	JWTSecret      string
	PolkaApiKey    string
	FileServerHits int
}

type DBConfig struct {
	DB *database.DB
}

type UserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Token       string `json:"token"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type ResBody struct {
	Body string `json:"body"`
}

type WebhookBody struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}
