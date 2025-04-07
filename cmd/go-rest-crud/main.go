package main

import (
	"log"
	"net/http"

	"github.com/RameezAhmad007/go-rest-crud/internal/handler"
)

func main() {

	// Set up HTTP routes
	mux := http.NewServeMux()
	handler.RegisterCardRoutes(mux)

	// Start server
	log.Println("Server starting on :8000...")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
