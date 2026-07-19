package handlers

import (
	"bytes"
	"fmt"
	"net/http"
)

// ErrorPageData contains the information displayed by error.html.
type ErrorPageData struct {
	StatusCode int
	Title      string
	Message    string
}

// renderError renders the custom error page and sends the correct HTTP status.
func renderError(
	w http.ResponseWriter,
	statusCode int,
	title string,
	message string,
) {
	data := ErrorPageData{
		StatusCode: statusCode,
		Title:      title,
		Message:    message,
	}

	var page bytes.Buffer

	err := Tmpl.ExecuteTemplate(
		&page,
		"error.html",
		data,
	)
	if err != nil {
		fmt.Println("Error Template Execution Error:", err)

		http.Error(
			w,
			http.StatusText(statusCode),
			statusCode,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(statusCode)

	_, err = page.WriteTo(w)
	if err != nil {
		fmt.Println("Error Response Write Error:", err)
	}
}
