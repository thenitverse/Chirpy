package main

import (
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserveHits atomic.Int32
	db            *database.Queries
	platform      string
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
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		return
	}
	err := cfg.db.DeleteUser(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}
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
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	platForm := os.Getenv("PLATFORM")

	myapiConfig := &apiConfig{
		db:       dbQueries,
		platform: platForm,
	}
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
	mux.HandleFunc("POST /api/users", myapiConfig.handlerUser)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}
