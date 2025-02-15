package main

import (
	"net/http"
)

// The routes function returns a ServeMux with the application routes defined.
// servermux is a request multiplexer that matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
func (app application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("GET /{$}", app.app_home)
	mux.HandleFunc("POST /zfs_create", app.zfs_vol_create)
	mux.HandleFunc("POST /pg_create", app.pg_create)
	mux.HandleFunc("POST /pg_destroy", app.pg_destroy)
	return mux
}
