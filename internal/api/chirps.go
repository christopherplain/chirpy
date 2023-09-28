package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/christopherplain/chirpy/internal/model"
	"github.com/go-chi/chi/v5"
)

func (cfg ApiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg ApiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	token, err := model.ValidateJWT(tokenString, cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user")
		return
	}
	userID, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user")
		return
	}

	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if userID != chirp.AuthorID {
		respondWithError(w, http.StatusForbidden, "Forbidden request")
	}

	err = cfg.DB.DeleteChirp(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (cfg ApiConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	var authorID *int
	authorIDQuery, err := strconv.Atoi(r.URL.Query().Get("author_id"))
	if err == nil {
		*authorID = authorIDQuery
	}

	sort := "asc"
	sortQuery := r.URL.Query().Get("sort")
	if sortQuery == "desc" {
		sort = "desc"
	}

	chirps, err := cfg.DB.GetChirps(authorID, sort)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *ApiConfig) HandlePostChirp(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	token, err := model.ValidateJWT(tokenString, cfg.JwtSecret)
	if err != nil {
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

	reqBody := model.Chirp{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	if err = model.ValidateChirp(reqBody.Body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	savedChirp, err := cfg.DB.CreateChirp(reqBody.Body, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, savedChirp)
}
