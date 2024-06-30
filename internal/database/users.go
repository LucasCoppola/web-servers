package database

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	newId := len(dbStructure.Users) + 1
	user := User{Id: newId, Email: email}

	if dbStructure.Users == nil {
		dbStructure.Users = make(map[int]User)
	}
	dbStructure.Users[newId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
