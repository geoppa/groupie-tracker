package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os" // Required to read environment variables

	"groupie-tracker/internal/handlers"
)

func main() {
	var err error

	// Parsing templates from the root execution path perspective
	handlers.Tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Error parsing templates:", err)
		return
	}

	// Static Asset Routing
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	// Dynamic URL Route Handling
	http.HandleFunc("/", handlers.ArtistsHandler)
	http.HandleFunc("/artist", handlers.ArtistDetailHandler)
	http.HandleFunc("/search", handlers.SearchHandler)

	// FIX: Dynamically fetch the port assigned by the hosting provider
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local development
	}

	fmt.Printf("Server starting on port %s...\n", port)

	// FIX: Listen on the assigned dynamic port
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
