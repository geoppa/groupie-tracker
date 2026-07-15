package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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

// ArtistDetailHandler processes requests for specific profiles (e.g., /artist?id=2)
func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the "id" parameter from the URL query string
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing artist ID parameter", http.StatusBadRequest)
		return
	}

	// 2. Convert the ID parameter from string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID format", http.StatusBadRequest)
		return
	}

	// 3. Retrieve target data components from the API package
	artist, err := api.FetchArtistByID(id)
	if err != nil {
		http.Error(w, "Artist data not found", http.StatusNotFound)
		return
	}

	relation, err := api.FetchRelation(id)
	if err != nil {
		http.Error(w, "Failed to fetch concert relations", http.StatusInternalServerError)
		return
	}

	// 4. Bundle data models inside an anonymous struct for the frontend template
	data := struct {
		Artist    interface{}
		Relations interface{}
	}{
		Artist:    artist,
		Relations: relation,
	}

	// 5. Inject bundled data into the specific details template
	err = Tmpl.ExecuteTemplate(w, "artist.html", data)
	if err != nil {
		fmt.Println("Template Execution Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
