package main

import (
	"log"
	"net/http"
)

// Create a new ServeMux
func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}

func main() {
	// Define the server and use the new ServeMux as the handler
	var mux = NewServeMux()
	var server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server at port :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}
