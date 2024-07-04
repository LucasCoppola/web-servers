package database

import (
	"errors"
	"time"
)

func (db *DB) StoreExpirationInSecs(userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, exists := dbStructure.Users[userId]
	if !exists {
		return errors.New("User not found")
	}

	expirationTimestamp := refreshTokenExpirationTime()
	user.RefreshTokenExpiresAt = expirationTimestamp

	dbStructure.Users[userId] = user

	return db.writeDB(dbStructure)
}

func (db *DB) StoreRefreshToken(refreshToken string, userId int, expiresAt int64) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, exists := dbStructure.Users[userId]
	if !exists {
		return errors.New("User not found")
	}

	user.RefreshToken = refreshToken
	user.RefreshTokenExpiresAt = expiresAt

	dbStructure.Users[userId] = user

	return db.writeDB(dbStructure)
}

func (db *DB) FindRefreshToken(refreshToken string) (bool, int64, int, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, 0, 0, err
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken {
			return true, user.RefreshTokenExpiresAt, user.Id, nil
		}
	}

	return false, 0, 0, errors.New("Couldn't find the refresh token")
}

func (db *DB) RevokeRefreshToken(userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, exists := dbStructure.Users[userId]

	if !exists {
		return errors.New("Couldn't find user")
	}

	user.RefreshToken = ""
	user.RefreshTokenExpiresAt = 0
	dbStructure.Users[userId] = user

	return db.writeDB(dbStructure)
}

func refreshTokenExpirationTime() int64 {
	duration := time.Hour * 24 * 60
	now := time.Now()
	expirationTime := now.Add(duration)
	expirationTimestamp := expirationTime.Unix()

	return expirationTimestamp
}
