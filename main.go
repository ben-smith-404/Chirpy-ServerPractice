package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ben-smith-404/Chirpy-ServerPractice/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		exitProgram(err)
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		isDev:          false,
	}
	if os.Getenv("PLATFORM") == "dev" {
		apiCfg.isDev = true
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handleReady)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleNewChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handleAddUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		exitProgram(err)
	}
}

func exitProgram(err error) {
	fmt.Println(err)
	os.Exit(1)
}
