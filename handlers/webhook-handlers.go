package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (dbCfg *DBConfig) WebhookHandler(w http.ResponseWriter, r *http.Request, polkaApiKey string) {
	header := r.Header.Get("Authorization")
	apiKey := strings.TrimPrefix(header, "ApiKey ")

	if apiKey == "" || apiKey != polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var webhook WebhookBody
	err := decoder.Decode(&webhook)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if webhook.Event == "user.upgraded" {
		code, err := dbCfg.DB.UpgradeUser(webhook.Data.UserId)
		if err != nil {
			respondWithError(w, code, err.Error())
			return
		}

		w.WriteHeader(code)
		return
	} else {
		w.WriteHeader(204)
		return
	}
}
