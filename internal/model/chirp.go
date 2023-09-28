package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	body = cleanBody(body, badWords)

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: authorID,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, ok := dbStructure.Chirps[id]; ok {
		delete(dbStructure.Chirps, id)
		return nil
	}

	return fmt.Errorf("Chirp ID %d not found", id)
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if chirp, ok := dbStructure.Chirps[id]; ok {
		return chirp, nil
	}

	return Chirp{}, fmt.Errorf("Chirp ID %d not found", id)
}

func (db *DB) GetChirps(authorID *int, order string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		if authorID != nil && chirp.AuthorID != *authorID {
			continue
		}
		chirps = append(chirps, chirp)
	}

	lessFunc := func(i, j int) bool { return chirps[i].ID < chirps[j].ID }
	if order == "desc" {
		lessFunc = func(i, j int) bool { return chirps[i].ID > chirps[j].ID }
	}
	sort.Slice(chirps, lessFunc)

	return chirps, nil
}

func ValidateChirp(body string) error {
	if len(body) > 140 {
		return errors.New("Chirp is too long")
	}
	return nil
}

func cleanBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	cleanedBody := strings.Join(words, " ")
	return cleanedBody
}
