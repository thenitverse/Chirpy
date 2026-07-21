package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserveHits atomic.Int32
}
type chirpRequest struct {
	Body string `json:"body"`
}
type validResponse struct {
	Valid bool `json:"valid"`
}
type errorResponse struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserveHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserveHits.Load()
	msg := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(msg))
}
func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserveHits.Store(0)
	w.WriteHeader(200)
}

func (cfg *apiConfig) validate_chirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	mychirpRequest := chirpRequest{}
	err := decoder.Decode(&mychirpRequest)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(mychirpRequest.Body) > 140 {
		respBody := errorResponse{Error: "Chirp is too long"}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return

	} else {
		cleaned := cleanChirp(mychirpRequest.Body)
		respBody := cleanedResponse{
			CleanedBody: cleaned,
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	}

}

func main() {
	myapiConfig := &apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", myapiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", myapiConfig.handleMetrics)
	mux.HandleFunc("POST /admin/reset", myapiConfig.handleReset)
	mux.HandleFunc("POST /api/validate_chirp", myapiConfig.validate_chirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}
