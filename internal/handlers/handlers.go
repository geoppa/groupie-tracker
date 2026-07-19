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
		renderError(
			w,
			http.StatusNotFound,
			"page_not_found",
			"Page Not Found",
			"The requested page could not be found.",
		)
		return
	}

	artists, err := api.FetchArtists()
	if err != nil {
		fmt.Println("Fetch Artists Error:", err)

		renderError(
			w,
			http.StatusInternalServerError,
			"artists_fetch_failed",
			"Internal Server Error",
			"The artist collection could not be loaded.",
		)
		return
	}

	pageData, err := newIndexPageData(artists)
	if err != nil {
		fmt.Println("Artist JSON Encoding Error:", err)

		renderError(
			w,
			http.StatusInternalServerError,
			"artist_data_failed",
			"Internal Server Error",
			"The artist data could not be prepared.",
		)
		return
	}

	err = Tmpl.ExecuteTemplate(
		w,
		"index.html",
		pageData,
	)
	if err != nil {
		fmt.Println("Template Execution Error:", err)

		/*
			The response may already have started here.

			For now, renderError is still our fallback.
			Later, during testing, we can improve successful-page
			rendering with a buffer as well.
		*/
		renderError(
			w,
			http.StatusInternalServerError,
			"template_failed",
			"Internal Server Error",
			"The page could not be prepared.",
		)
		return
	}
}

// ArtistDetailHandler processes requests for a specific artist.
func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		renderError(
			w,
			http.StatusBadRequest,
			"missing_artist_id",
			"Bad Request",
			"No artist was selected.",
		)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(
			w,
			http.StatusBadRequest,
			"invalid_artist_id",
			"Bad Request",
			"The artist identifier is invalid.",
		)
		return
	}

	if id <= 0 {
		renderError(
			w,
			http.StatusBadRequest,
			"invalid_artist_id",
			"Bad Request",
			"The artist identifier is invalid.",
		)
		return
	}

	artist, err := api.FetchArtistByID(id)
	if err != nil {
		fmt.Println("Fetch Artist Error:", err)

		renderError(
			w,
			http.StatusNotFound,
			"artist_not_found",
			"Artist Not Found",
			"The requested artist could not be found.",
		)
		return
	}

	relation, err := api.FetchRelation(id)
	if err != nil {
		fmt.Println("Fetch Relation Error:", err)

		renderError(
			w,
			http.StatusInternalServerError,
			"relation_fetch_failed",
			"Internal Server Error",
			"The concert information could not be loaded.",
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

	err = Tmpl.ExecuteTemplate(
		w,
		"artist.html",
		data,
	)
	if err != nil {
		fmt.Println("Template Execution Error:", err)

		renderError(
			w,
			http.StatusInternalServerError,
			"template_failed",
			"Internal Server Error",
			"The page could not be prepared.",
		)
		return
	}
}
