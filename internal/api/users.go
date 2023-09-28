package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/christopherplain/chirpy/internal/model"
)

type UserReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRespBody struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (cfg ApiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	reqBody := UserReqBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	user, err := cfg.DB.AuthenticateUser(reqBody.Email, reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respBody := UserRespBody{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	token, err := model.GenerateAccessToken(cfg.JwtSecret, user.ID)
	if err == nil {
		respBody.Token = token
	}

	token, err = model.GenerateRefreshToken(cfg.JwtSecret, user.ID)
	if err == nil {
		respBody.RefreshToken = token
	}

	respondWithJSON(w, http.StatusOK, respBody)
}

func (cfg ApiConfig) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	reqBody := UserReqBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	savedUser, err := cfg.DB.CreateUser(reqBody.Email, reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respBody := UserRespBody{
		ID:          savedUser.ID,
		Email:       savedUser.Email,
		IsChirpyRed: savedUser.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, respBody)
}

func (cfg ApiConfig) HandlePutUser(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	token, err := model.ValidateJWT(tokenString, cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer == "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user")
		return
	}

	reqBody := UserReqBody{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	user, err := cfg.DB.UpdateUser(id, reqBody.Email, reqBody.Password, nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respBody := UserRespBody{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}
