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
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK")) //converting string "OK" to a slice of bytes
	})

	srv := &http.Server{ //type http.Server
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port:%s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
	//anything after this line will not work as the previous line blocks

}