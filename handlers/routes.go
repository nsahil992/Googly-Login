package handlers

import (
	"net/http"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	// API endpoints
	mux.HandleFunc("/api/register", RegisterHandler)
	mux.HandleFunc("/api/login", LoginHandler)

	return mux
}
