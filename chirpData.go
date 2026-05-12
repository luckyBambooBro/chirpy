package main
import (
	"encoding/json"
	"net/http"
	"log"
)

var minChirpLength = 140

type chirpData struct {
	Content string `json:"body"`
}

type chirpError struct {
	ChirpError string `json:"error"`
}

type chirpValid struct {
	ChirpValid bool `json:"valid"`
}

func respondWithError(w http.ResponseWriter, code int, errorMsg string) {
	//encode errorMsg into JSON here
	chirpErrorMsg := &chirpError{
		ChirpError: errorMsg,
	}
	data, err := json.Marshal(chirpErrorMsg) 
	if err != nil {
		log.Printf("unable to encode chirp error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON (w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("unable to encode valid chirp response: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	//decode request
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	chirp := &chirpData{}
	if err := decoder.Decode(chirp); err != nil {
		log.Printf("Error decoding request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "unable to decode request")
		return
	}
	
	//handle request depending on length
	if len(chirp.Content) > minChirpLength {
		respondWithError(w , http.StatusBadRequest, "Chirp is too long")
		return //is this return necessary???
	} else {
		respondWithJSON(w, http.StatusOK, chirpValid{
			ChirpValid: true,
		})
	}

}