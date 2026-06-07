package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/luckyBambooBro/chirpy/internal/database"
)

func main() {
	const (
		filePathRoot = "."
		port         = "8080"
	)

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("%v must be set", dbURL)
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT secret must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		db:       dbQueries,
		platform: platform,
		jwtSecret: jwtSecret,
	}
	mux := http.NewServeMux() //type: *http.ServeMux
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGetOne)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerUserRefresh)

	srv := &http.Server{ //type http.Server
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port:%s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
	//anything after this line will not work as the previous line blocks
}
