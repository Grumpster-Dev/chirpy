package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Grumpster-Dev/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// func (cfg *apiConfig) getFileserverHits() int32 {
// 	return cfg.fileserverHits.Load()
// }

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func fileServerHandler() http.Handler {
	return http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	x := cfg.fileserverHits.Load()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, x)

	w.Write([]byte(body))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Reset endpoint is only available in dev environment")
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to delete users")
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits counter reset to 0"))

}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type chirpResponse struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    string    `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	post := chirpRequest{}
	err := decoder.Decode(&post)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	// post is a struct with data populated successfully
	if len(post.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleaned := cleanedBody(post.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: post.UserID})

	if err != nil {
		respondWithError(w, 500, "Could not create chirp")
		return
	}
	respondWithJSON(w, 201, chirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := fmt.Sprintf(`{"error":"%s"}`, msg)
	w.Write([]byte(response))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func cleanedBody(chirp string) string {
	temp := strings.Split(chirp, " ")
	for i, word := range temp {
		lowered := strings.ToLower(word)
		if lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
			temp[i] = "****"
		}
	}
	return strings.Join(temp, " ")

}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	req := userRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, 500, "Could not create user")
		return
	}

	resp := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, 201, resp)
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(db)

	cfg := &apiConfig{
		db:       dbQueries,
		platform: platform}
	mux := http.NewServeMux()
	handler := fileServerHandler()
	wrapped := cfg.middlewareMetricsInc(handler)

	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.Handle("/app/", wrapped)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
