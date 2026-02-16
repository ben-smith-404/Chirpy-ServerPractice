package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/ben-smith-404/Chirpy-ServerPractice/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	isDev          bool
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleHits(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")
	htmlString := `<html>
	<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
	</body>
	</html>`
	fmt.Fprintf(w, htmlString, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if !cfg.isDev {
		log.Print("Endpoint only available in development")
		respondWithError(w, 403, "something went wrong")
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Error resetting user: %s", err)
		respondWithError(w, 500, "something went wrong")
		return
	}
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
}
