package main

import (
	"encoding/json"
	"net/http"
	"text/template"
)

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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type ResBody struct {
		Body string `json:"body"`
	}
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
	validRes := SuccessRes{Body: nonProfaneMsg}

	respondWithJSON(w, 200, validRes)
}