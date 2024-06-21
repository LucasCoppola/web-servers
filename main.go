package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	filserverHits int
}

func main() {
	const PORT = "8080"
	mux := http.NewServeMux()

	var apiCfg apiConfig

	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.numOfReqsHandler)
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetNumOfReqsHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
