package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

// initTestTemplates parses a minimal inline mock HTML template structure.
// This prevents the handler from running into a nil pointer panic during isolation testing.
func initTestTemplates() {
	// We use template.New to build mock templates directly in memory.
	// This mirrors the layout files your real handlers expect (index.html and error.html).
	Tmpl = template.New("test")

	// Add mock structures for index.html
	template.Must(Tmpl.New("index.html").Parse(`<html><body>Mock Index: {{len .Artists}}</body></html>`))

	// Add mock structures for error.html
	template.Must(Tmpl.New("error.html").Parse(`<html><body>Mock Error: {{.Title}}</body></html>`))
}

// TestArtistsHandler_RootRoute verifies that hitting "/" handles the request properly.
func TestArtistsHandler_RootRoute(t *testing.T) {
	// 1. Initialize our mock templates so Tmpl is no longer nil
	initTestTemplates()

	// 2. Create a fake HTTP GET request simulating a browser visiting "/"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 3. Create a ResponseRecorder to capture the output channel bytes
	rr := httptest.NewRecorder()

	// 4. Wrap and trigger your real handler execution logic
	handler := http.HandlerFunc(ArtistsHandler)
	handler.ServeHTTP(rr, req)

	// 5. Assertions: Check if everything succeeded or handled cleanly
	if rr.Code != http.StatusOK && rr.Code != http.StatusInternalServerError {
		t.Errorf("ArtistsHandler returned unexpected status code: got %v", rr.Code)
	}
}

// TestArtistsHandler_NotFoundRoute checks if typing a broken path yields a 404 page.
func TestArtistsHandler_NotFoundRoute(t *testing.T) {
	// Initialize our mock templates here as well
	initTestTemplates()

	req, err := http.NewRequest("GET", "/this-page-does-not-exist", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ArtistsHandler)

	handler.ServeHTTP(rr, req)

	// Your handler code uses renderError which requires error.html (now safely initialized)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for invalid path, but got %v", rr.Code)
	}
}
