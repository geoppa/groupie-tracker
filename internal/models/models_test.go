package models

import (
	"encoding/json"
	"testing"
)

// TestArtistJSONMapping checks if the Artist struct correctly unmarshals JSON keys.
func TestArtistJSONMapping(t *testing.T) {
	jsonData := `{
		"id": 1,
		"name": "Queen",
		"members": ["Freddie Mercury", "Brian May"],
		"creationDate": 1970,
		"firstAlbum": "13-07-1973"
	}`

	var artist Artist
	err := json.Unmarshal([]byte(jsonData), &artist)
	if err != nil {
		t.Fatalf("Failed to unmarshal Artist JSON: %v", err)
	}

	if artist.Name != "Queen" {
		t.Errorf("Expected name 'Queen', got '%s'", artist.Name)
	}

	if artist.CreationDate != 1970 {
		t.Errorf("Expected creation date 1970, got %d", artist.CreationDate)
	}
}

// TestRelationMap checks if our custom nested map parses the API response correctly.
func TestRelationMap(t *testing.T) {
	jsonData := `{
		"id": 1,
		"datesLocations": {
			"london-uk": ["12-12-2019", "13-12-2019"],
			"paris-france": ["05-05-2020"]
		}
	}`

	var relation Relation
	err := json.Unmarshal([]byte(jsonData), &relation)
	if err != nil {
		t.Fatalf("Failed to unmarshal Relation JSON: %v", err)
	}

	dates, exists := relation.DatesLocations["london-uk"]
	if !exists {
		t.Fatal("Expected key 'london-uk' to exist in map")
	}

	if len(dates) != 2 {
		t.Errorf("Expected 2 dates for London, got %d", len(dates))
	}
}
