package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

// SearchHandler evaluates queries against names, members, locations, and dates
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get and sanitize search query parameters
	query := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("q")))
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 2. Fetch Artists and Locations
	allArtists, err := api.FetchArtists()
	if err != nil {
		fmt.Println("Search - Fetch Artists Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	locIndex, err := api.FetchLocations()
	if err != nil {
		fmt.Println("Search - Fetch Locations Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Map locations by Artist ID for quick lookups during iteration
	locationMap := make(map[int][]string)
	for _, node := range locIndex.Index {
		locationMap[node.ID] = node.Locations
	}

	var filteredArtists []models.Artist

	// 3. Deep search algorithm parsing
	for _, artist := range allArtists {
		// Rule A: Match against Band/Artist Name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		// Rule B: Match against Creation Date
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		// Rule C: Match against First Album Release Date String
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		// Rule D: Match against Individual Band Members
		memberMatch := false
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				memberMatch = true
				break
			}
		}
		if memberMatch {
			filteredArtists = append(filteredArtists, artist)
			continue
		}

		// Rule E: Match against Concert Locations
		locationMatch := false
		if locs, exists := locationMap[artist.ID]; exists {
			for _, loc := range locs {
				// Cleans up location formatting strings (e.g., "london-uk" -> "london uk")
				cleanLoc := strings.ReplaceAll(loc, "-", " ")
				cleanLoc = strings.ReplaceAll(cleanLoc, "_", " ")
				if strings.Contains(strings.ToLower(cleanLoc), query) {
					locationMatch = true
					break
				}
			}
		}
		if locationMatch {
			filteredArtists = append(filteredArtists, artist)
			continue
		}
	}

	// 4. Render results using your global template variable Tmpl
	err = Tmpl.ExecuteTemplate(w, "index.html", filteredArtists)
	if err != nil {
		fmt.Println("Search Template Execution Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
