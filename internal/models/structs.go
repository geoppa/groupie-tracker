package models

// Artist represents the data structure for a single artist from the API
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocationsURL string   `json:"locations"`
	DatesURL     string   `json:"concertDates"`
	RelationsURL string   `json:"relations"`
}

// Relation represents the outer wrapper structure returned by the /relation endpoint
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"` // Maps city names to their respective event dates
}

// LocationIndex represents the top-level array container returned by the /locations endpoint
type LocationIndex struct {
	Index []LocationNode `json:"index"`
}

// LocationNode represents the specific locations assigned to a single band ID
type LocationNode struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}
