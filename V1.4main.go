package main

import (
	"fmt"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func fileServerHandler() http.Handler {
	return http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthzHandler)
	mux.Handle("/app/", fileServerHandler())

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
