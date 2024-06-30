package main

import (
	"github.com/LucasCoppola/web-server/internal/database"
	"log"
	"net/http"
)

type apiConfig struct {
	filserverHits int
}

type dbConfig struct {
	DB *database.DB
}

func main() {
	const PORT = "8080"
	mux := http.NewServeMux()

	db, err := database.NewDB("internal/database/db.json")

	if err != nil {
		log.Fatal(err)
	}

	dbCfg := &dbConfig{
		DB: db,
	}
	apiCfg := &apiConfig{
		filserverHits: 0,
	}

	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.numOfReqsHandler)
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetNumOfReqsHandler)
	mux.HandleFunc("GET /api/chirps", dbCfg.getChirpHandler)
	mux.HandleFunc("POST /api/chirps", dbCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps/{id}", dbCfg.getSingleChirpHandler)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
