package main

import (
	"log"
	"net/http"
	"text/template"
)

type apiConfig struct {
	filserverHits int
}

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.filserverHits++
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	const PORT = "8080"
	mux := http.NewServeMux()

	var apiCfg apiConfig

	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handler)
	mux.HandleFunc("/admin/metrics", apiCfg.numOfReqsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetNumOfReqsHandler)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
