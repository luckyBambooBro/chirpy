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
	mux := http.NewServeMux() //type: *http.ServeMux
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))
	
	srv := &http.Server{ //type http.Server
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port:%s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
	//up to here
}

