package main

import (
	"encoding/json"
	"net/http"
)

// Artist represents the data structure for a single artist from the API
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocationsURL string   `json:"locations"`    // URL to fetch locations
	DatesURL     string   `json:"concertDates"` // URL to fetch dates
	RelationsURL string   `json:"relations"`    // URL to fetch relations
}

// fetchArtists makes an HTTP GET request to retrieve the full list of artists
func fetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://herokuapp.com") // HTTP GET request
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	// Decode the JSON response body directly into the artists slice
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}
