package main

import (
	"fmt"
	"html/template"
	"net/http"

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

	fmt.Println("Server starting on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
