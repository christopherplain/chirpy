package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PolkaReqBody struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func (cfg ApiConfig) HandlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	authHeader := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authHeader) < 2 || authHeader[1] != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Not authorized")
	}

	reqBody := PolkaReqBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	if reqBody.Event == "user.upgraded" {
		isChirpyRed := true
		_, err = cfg.DB.UpdateUser(reqBody.Data.UserID, "", "", &isChirpyRed)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
