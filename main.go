package main

import (
	"github.com/LucasCoppola/web-server/handlers"
	"github.com/LucasCoppola/web-server/internal/database"
	"log"
	"net/http"
)

func main() {
	const PORT = "8080"
	mux := http.NewServeMux()

	db, err := database.NewDB("internal/database/db.json")

	if err != nil {
		log.Fatal(err)
	}

	dbCfg := &handlers.DBConfig{
		DB: db,
	}

	apiCfg := &handlers.ApiConfig{
		FileServerHits: 0,
	}

	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.NumOfReqsHandler)
	mux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.ResetNumOfReqsHandler)

	mux.HandleFunc("GET /api/chirps", dbCfg.GetChirpHandler)
	mux.HandleFunc("POST /api/chirps", dbCfg.CreateChirpHandler)
	mux.HandleFunc("GET /api/chirps/{id}", dbCfg.GetSingleChirpHandler)

	mux.HandleFunc("POST /api/users", dbCfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", dbCfg.LoginUserHandler)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
