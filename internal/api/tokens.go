package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/christopherplain/chirpy/internal/model"
)

type RefreshRespBody struct {
	Token string `json:"token"`
}

func (cfg ApiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	token, err := model.ValidateJWT(tokenString, cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}
	isRevoked, err := cfg.DB.IsTokenRevoked(tokenString)
	if err != nil || isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newToken, err := model.GenerateAccessToken(cfg.JwtSecret, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respBody := RefreshRespBody{
		Token: newToken,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}

func (cfg ApiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	token, err := model.ValidateJWT(tokenString, cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}
	isRevoked, err := cfg.DB.IsTokenRevoked(tokenString)
	if err != nil || isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}

	err = cfg.DB.RevokeToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
