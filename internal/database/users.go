package database

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
