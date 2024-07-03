package database

import (
	"errors"
)

func (db *DB) CreateUser(email string, hashedPassword []byte) (UserResponse, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return UserResponse{}, err
	}

	newId := len(dbStructure.Users) + 1
	user := User{Id: newId, Email: email, Password: string(hashedPassword)}

	if dbStructure.Users == nil {
		dbStructure.Users = make(map[int]User)
	}
	dbStructure.Users[newId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{Id: user.Id, Email: user.Email}, nil
}

func (db *DB) UpdateUser(userId int, email string, hashedPassword []byte) (UserResponse, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return UserResponse{}, err
	}

	user, exists, err := db.FindUserById(userId)

	if err != nil {
		return UserResponse{}, err
	}

	if !exists {
		return UserResponse{}, errors.New("User doesn't exists")
	}

	user.Email = email
	user.Password = string(hashedPassword)

	dbStructure.Users[user.Id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{Id: user.Id, Email: user.Email}, nil
}

func (db *DB) FindUserById(userId int) (User, bool, error) {
	DBStructure, err := db.loadDB()

	if err != nil {
		return User{}, false, err
	}

	user, exists := DBStructure.Users[userId]

	return user, exists, nil
}

func (db *DB) FindUserByEmail(email string) (User, bool, error) {
	DBStructure, err := db.loadDB()

	if err != nil {
		return User{}, false, err
	}

	for _, user := range DBStructure.Users {
		if user.Email == email {
			return user, true, nil
		}
	}

	return User{}, false, nil
}

func (db *DB) StoreToken(userId int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, exists := dbStructure.Users[userId]
	if !exists {
		return errors.New("User not found")
	}

	user.Token = token
	dbStructure.Users[userId] = user

	return db.writeDB(dbStructure)
}
