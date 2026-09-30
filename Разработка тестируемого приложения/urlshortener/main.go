package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	shortener := NewURLShortener()

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler(shortener))
	mux.HandleFunc("/", redirectHandler(shortener))

	log.Println("URL shortener is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func shortenHandler(shortener *URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var request shortenRequest
		if err := decoder.Decode(&request); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		shortID, err := shortener.Shorten(request.URL)
		if err != nil {
			if errors.Is(err, ErrInvalidURL) {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}

			writeJSONError(w, http.StatusInternalServerError, "could not shorten URL")
			return
		}

		response := shortenResponse{
			ShortURL:    shortID,
			OriginalURL: request.URL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("cannot encode response: %v", err)
		}
	}
}

func redirectHandler(shortener *URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		shortID := strings.TrimPrefix(r.URL.Path, "/")

		if shortID == "" || strings.Contains(shortID, "/") {
			writeJSONError(w, http.StatusNotFound, "short URL not found")
			return
		}

		originalURL, err := shortener.GetOriginal(shortID)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "short URL not found")
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorResponse{Error: message}); err != nil {
		log.Printf("cannot encode error response: %v", err)
	}
}
