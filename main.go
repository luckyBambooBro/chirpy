package main

import (
	"log"
	"net/http"
)

func main() {
	const (
		filePathRoot = "."
		port = "8080"
	)
	apiCfg := apiConfig{}
	mux := http.NewServeMux() //type: *http.ServeMux
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.apiHandler)
	mux.HandleFunc("/reset", apiCfg.reset)

	srv := &http.Server{ //type http.Server
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port:%s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
	//anything after this line will not work as the previous line blocks
}



