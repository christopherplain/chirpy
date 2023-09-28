package model

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (db *DB) AuthenticateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	errMsg := "Invalid email address or password"
	id, ok := dbStructure.UsersEmailToID[email]
	if !ok {
		return User{}, errors.New(errMsg)
	}
	user := dbStructure.Users[id]

	if !doPasswordsMatch(user.Password, password) {
		return User{}, errors.New(errMsg)
	}

	return user, nil
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if id, ok := dbStructure.UsersEmailToID[email]; ok {
		return dbStructure.Users[id], nil
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %s\n", err)
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:          id,
		Email:       email,
		Password:    hashedPassword,
		IsChirpyRed: false,
	}
	dbStructure.Users[id] = user
	dbStructure.UsersEmailToID[email] = id

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(id int, email string, password string, isChirpyRed *bool) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, fmt.Errorf("unable to fetch user with ID %d", id)
	}
	updatedUser := User{
		ID:          id,
		Email:       user.Email,
		Password:    user.Password,
		IsChirpyRed: user.IsChirpyRed,
	}

	if email != "" {
		updatedUser.Email = email
		dbStructure.UsersEmailToID[email] = id
	}

	if password != "" {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			log.Printf("Error hashing password: %s\n", err)
			return User{}, err
		}
		updatedUser.Password = hashedPassword
	}

	if isChirpyRed != nil {
		updatedUser.IsChirpyRed = *isChirpyRed
	}

	dbStructure.Users[id] = updatedUser
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func doPasswordsMatch(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashedPassword), err
}
