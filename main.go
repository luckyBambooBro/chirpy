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
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("%v must be set", dbURL)
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	dbQueries := database.New(dbConn)

	const (
		filePathRoot = "."
		port         = "8080"
	)
	apiCfg := apiConfig{
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}
	mux := http.NewServeMux() //type: *http.ServeMux
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	srv := &http.Server{ //type http.Server
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port:%s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
	//anything after this line will not work as the previous line blocks
}

func returnErrorMsg(w http.ResponseWriter, code int, errorMsg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	if _, err := w.Write([]byte(errorMsg)); err != nil {
		log.Println(err)
	}
}
