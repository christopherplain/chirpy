package model

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (db *DB) IsTokenRevoked(token string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}
	if _, ok := dbStructure.RevokedTokens[token]; ok {
		return true, nil
	}
	return false, nil
}

func (db *DB) RevokeToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	dbStructure.RevokedTokens[token] = time.Now().String()
	return db.writeDB(dbStructure)
}

func GenerateAccessToken(secret string, id int) (string, error) {
	issuer := "chirpy-access"
	expires := time.Duration(1 * time.Hour)
	return generateJWT(secret, issuer, expires, id)
}

func GenerateRefreshToken(secret string, id int) (string, error) {
	issuer := "chirpy-refresh"
	expires := time.Duration(1440 * time.Hour)
	return generateJWT(secret, issuer, expires, id)
}

func ValidateJWT(token string, secret string) (*jwt.Token, error) {
	claims := &jwt.RegisteredClaims{}
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func generateJWT(secret string, issuer string, expires time.Duration, id int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
		Subject:   strconv.Itoa(id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
