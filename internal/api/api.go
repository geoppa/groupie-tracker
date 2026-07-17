package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"groupie-tracker/internal/models"
)

// BaseAPI Hardcodes the root URL of the project's external data source
const BaseAPI = "https://groupietrackers.herokuapp.com/api"

// FetchArtists connects directly to the sub-endpoint
func FetchArtists() ([]models.Artist, error) {
	// Sends a request to retrieve the full list of bands
	resp, err := http.Get(BaseAPI + "/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}

// FetchArtistByID fetches a single artist profile by their unique ID
func FetchArtistByID(id int) (models.Artist, error) {
	// Sprintf builds the exact web address by appending strings and an integer together
	resp, err := http.Get(fmt.Sprintf("%s/artists/%d", BaseAPI, id))
	if err != nil {
		return models.Artist{}, err
	}
	defer resp.Body.Close()

	var artist models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artist)
	return artist, err
}

// FetchRelation fetches combined concert locations and dates for a specific artist
func FetchRelation(id int) (models.Relation, error) {
	resp, err := http.Get(fmt.Sprintf("%s/relation/%d", BaseAPI, id))
	if err != nil {
		return models.Relation{}, err
	}
	defer resp.Body.Close()

	var relation models.Relation
	err = json.NewDecoder(resp.Body).Decode(&relation)
	return relation, err
}

// FetchLocations retrieves the raw location index data for all bands
func FetchLocations() (models.LocationIndex, error) {
	resp, err := http.Get(BaseAPI + "/locations")
	if err != nil {
		return models.LocationIndex{}, err
	}
	defer resp.Body.Close()

	var locIndex models.LocationIndex
	err = json.NewDecoder(resp.Body).Decode(&locIndex)
	return locIndex, err
}
