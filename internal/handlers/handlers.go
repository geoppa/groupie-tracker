package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"groupie-tracker/internal/api"
)

// Tmpl holds the pre-parsed templates globally
var Tmpl *template.Template

// ArtistsHandler handles the main domain route
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	artists, err := api.FetchArtists()
	if err != nil {
		fmt.Println("Fetch Artists Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = Tmpl.ExecuteTemplate(w, "index.html", artists)
	if err != nil {
		fmt.Println("Template Execution Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
