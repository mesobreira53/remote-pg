package main

import (
	"log"
	"net/http"
)

// The main function is the entry point for the application
func main() {

	// Create a new application
	app := newApplication("zpool")

	app.logger.Info("Starting server on :8080")
	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
