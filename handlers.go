package main

import (
	"encoding/json"
	"net/http"
	"text/template"
)

type ResBody struct {
	Body string `json:"body"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (apiCfg *apiConfig) numOfReqsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.New("admin").Parse(`
    <html>
      <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited {{.Hits}} times!</p>
      </body>
    </html>
    `))

	tmpl.Execute(w, struct{ Hits int }{Hits: apiCfg.filserverHits})
}

func (apiCfg *apiConfig) resetNumOfReqsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	apiCfg.filserverHits = 0
}

func (dbCfg *dbConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type SuccessRes struct {
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	body := ResBody{}
	err := decoder.Decode(&body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	nonProfaneMsg := removeProfanity(body.Body)

	chirp, err := dbCfg.DB.CreateChirp(nonProfaneMsg)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, 201, chirp)
}

func (dbCfg *dbConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := dbCfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, 200, chirps)
}

func (dbCfg *dbConfig) getSingleChirpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	chirp, err := dbCfg.DB.GetSingleChirp(id)

	if err != nil {
		if err.Error() == "chirp not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(w, 200, chirp)
}
