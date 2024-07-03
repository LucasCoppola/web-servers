package main

import (
	"github.com/LucasCoppola/web-server/handlers"
	"github.com/LucasCoppola/web-server/internal/database"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	const PORT = "8080"
	mux := http.NewServeMux()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewDB("internal/database/db.json")

	if err != nil {
		log.Fatal(err)
	}

	dbCfg := &handlers.DBConfig{
		DB: db,
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	apiCfg := &handlers.ApiConfig{
		JWTSecret:      jwtSecret,
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
	mux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, r *http.Request) {
		dbCfg.UpdateUserHandler(w, r, apiCfg.JWTSecret)
	})

	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		dbCfg.LoginHandler(w, r, apiCfg.JWTSecret)
	})

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
