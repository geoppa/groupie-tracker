package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

// SearchHandler evaluates queries against artist information.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(
		strings.ToLower(r.URL.Query().Get("q")),
	)

	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	allArtists, err := api.FetchArtists()
	if err != nil {
		fmt.Println("Search - Fetch Artists Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	locationIndex, err := api.FetchLocations()
	if err != nil {
		fmt.Println("Search - Fetch Locations Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	locationMap := make(map[int][]string)

	for _, locationNode := range locationIndex.Index {
		locationMap[locationNode.ID] = locationNode.Locations
	}

	var filteredArtists []models.Artist

	for _, artist := range allArtists {
		if strings.Contains(
			strings.ToLower(artist.Name),
			query,
		) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		if strings.Contains(
			strconv.Itoa(artist.CreationDate),
			query,
		) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		if strings.Contains(
			strings.ToLower(artist.FirstAlbum),
			query,
		) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		memberMatch := false

		for _, member := range artist.Members {
			if strings.Contains(
				strings.ToLower(member),
				query,
			) {
				memberMatch = true
				break
			}
		}

		if memberMatch {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		locationMatch := false

		locations, exists := locationMap[artist.ID]
		if exists {
			for _, location := range locations {
				cleanLocation := strings.ReplaceAll(
					location,
					"-",
					" ",
				)

				cleanLocation = strings.ReplaceAll(
					cleanLocation,
					"_",
					" ",
				)

				if strings.Contains(
					strings.ToLower(cleanLocation),
					query,
				) {
					locationMatch = true
					break
				}
			}
		}

		if locationMatch {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	pageData, err := newIndexPageData(filteredArtists)
	if err != nil {
		fmt.Println("Search - Artist JSON Encoding Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	err = Tmpl.ExecuteTemplate(w, "index.html", pageData)
	if err != nil {
		fmt.Println("Search Template Execution Error:", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}
}
