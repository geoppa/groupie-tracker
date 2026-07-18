package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

// Tmpl holds the pre-parsed templates globally.
var Tmpl *template.Template

// IndexPageData contains the data required by index.html.
type IndexPageData struct {
	Artists     []models.Artist
	ArtistsJSON template.JS
}

// newIndexPageData prepares the artist data for index.html
// and for the JavaScript autocomplete component.
func newIndexPageData(artists []models.Artist) (IndexPageData, error) {
	artistsJSON, err := json.Marshal(artists)
	if err != nil {
		return IndexPageData{}, err
	}

	return IndexPageData{
		Artists:     artists,
		ArtistsJSON: template.JS(artistsJSON),
	}, nil
}

// ArtistsHandler handles the main page route.
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	artists, err := api.FetchArtists()
	if err != nil {
		fmt.Println("Fetch Artists Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	pageData, err := newIndexPageData(artists)
	if err != nil {
		fmt.Println("Artist JSON Encoding Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	err = Tmpl.ExecuteTemplate(w, "index.html", pageData)
	if err != nil {
		fmt.Println("Template Execution Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}
}

// ArtistDetailHandler processes requests for a specific artist.
func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(
			w,
			"Missing artist ID parameter",
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(
			w,
			"Invalid artist ID format",
			http.StatusBadRequest,
		)
		return
	}

	artist, err := api.FetchArtistByID(id)
	if err != nil {
		http.Error(
			w,
			"Artist data not found",
			http.StatusNotFound,
		)
		return
	}

	relation, err := api.FetchRelation(id)
	if err != nil {
		http.Error(
			w,
			"Failed to fetch concert relations",
			http.StatusInternalServerError,
		)
		return
	}

	data := struct {
		Artist    models.Artist
		Relations models.Relation
	}{
		Artist:    artist,
		Relations: relation,
	}

	err = Tmpl.ExecuteTemplate(w, "artist.html", data)
	if err != nil {
		fmt.Println("Template Execution Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}
}
