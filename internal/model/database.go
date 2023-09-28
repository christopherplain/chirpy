package model

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type dbStructure struct {
	Chirps         map[int]Chirp     `json:"chirps"`
	Users          map[int]User      `json:"users"`
	UsersEmailToID map[string]int    `json:"users_email_to_id"`
	RevokedTokens  map[string]string `json:"revoked_tokens"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func ResetDB(path string) error {
	return os.Remove(path)
}

func (db *DB) createDB() error {
	dbStructure := dbStructure{
		Chirps:         map[int]Chirp{},
		Users:          map[int]User{},
		UsersEmailToID: map[string]int{},
		RevokedTokens:  map[string]string{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (dbStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := dbStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}

	err = json.Unmarshal(data, &dbStructure)

	return dbStructure, err
}

func (db *DB) writeDB(dbStructure dbStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0600)
}
