package handlers

import (
	"github.com/LucasCoppola/web-server/internal/database"
	"net/http"
	"text/template"
)

type ApiConfig struct {
	JWTSecret      string
	FileServerHits int
}

type DBConfig struct {
	DB *database.DB
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (apiCfg *ApiConfig) NumOfReqsHandler(w http.ResponseWriter, r *http.Request) {
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

	tmpl.Execute(w, struct{ Hits int }{Hits: apiCfg.FileServerHits})
}

func (apiCfg *ApiConfig) ResetNumOfReqsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	apiCfg.FileServerHits = 0
}
