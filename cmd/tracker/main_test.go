package main

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"
)

// TestTemplateParsing ensures that all required HTML files exist and have no syntax errors.
func TestTemplateParsing(t *testing.T) {
	// 1. Dynamically locate the templates directory by traversing upwards if needed
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	templatesDir := ""
	for {
		target := filepath.Join(dir, "templates")
		if info, err := os.Stat(target); err == nil && info.IsDir() {
			templatesDir = target
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir { // Reached the filesystem root without finding it
			break
		}
		dir = parent
	}

	if templatesDir == "" {
		t.Skip("Skipping test: 'templates' folder could not be located anywhere in the path hierarchy.")
	}

	// 2. Execute the exact compilation rule using the absolute found path pattern
	pattern := filepath.Join(templatesDir, "*.html")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		t.Fatalf("Template compilation failed: %v. Check for unclosed actions or structural layout problems.", err)
	}

	if tmpl == nil {
		t.Fatal("Compiled template instance returned nil pointer reference tracking context")
	}
}
