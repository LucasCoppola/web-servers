package handlers

import (
	"github.com/LucasCoppola/web-server/internal/database"
)

type ApiConfig struct {
	JWTSecret      string
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
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type ResBody struct {
	Body string `json:"body"`
}
