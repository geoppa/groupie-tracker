package api

import (
	"encoding/json"
	"net/http"

	"groupie-tracker/internal/models"
)

// BaseAPI is the exact URL specified for the project
const BaseAPI = "https://groupietrackers.herokuapp.com/api"

// FetchArtists connects directly to the sub-endpoint
func FetchArtists() ([]models.Artist, error) {
	// We concatenate BaseAPI + "/artists" to get the exact target URL
	resp, err := http.Get(BaseAPI + "/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}
