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

	http.HandleFunc("/", handlers.ArtistsHandler)

	fmt.Println("Server starting on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
