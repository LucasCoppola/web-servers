package database

import (
	"sync"
)

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

type User struct {
	Id                    int    `json:"id"`
	Email                 string `json:"email"`
	Password              string `json:"password"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
	IsChirpyRed           bool   `json:"is_chirpy_red"`
}

type UserResponse struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}
