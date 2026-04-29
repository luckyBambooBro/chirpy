package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux() //type: *http.ServeMux
	s := http.Server{ //type http.Server
		Addr: ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
	mux.Handle("/", http.FileServer(http.Dir(".")))

}

