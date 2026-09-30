package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantError  bool
	}{
		{
			name:       "successful shortening",
			method:     http.MethodPost,
			body:       `{"url":"https://example.com/very/long/path"}`,
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name:       "invalid JSON",
			method:     http.MethodPost,
			body:       `{"url":`,
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "invalid URL",
			method:     http.MethodPost,
			body:       `{"url":"not-a-url"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "unsupported method",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  true,
		},
		{
			name:       "unknown JSON field",
			method:     http.MethodPost,
			body:       `{"url":"https://example.com","extra":"field"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := NewURLShortener()
			handler := shortenHandler(shortener)

			request := httptest.NewRequest(
				tt.method,
				"/shorten",
				bytes.NewBufferString(tt.body),
			)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != tt.wantStatus {
				t.Fatalf(
					"status = %d, want %d; body = %s",
					response.StatusCode,
					tt.wantStatus,
					recorder.Body.String(),
				)
			}

			if tt.wantError {
				var errorBody errorResponse

				if err := json.NewDecoder(response.Body).Decode(&errorBody); err != nil {
					t.Fatalf("cannot decode error response: %v", err)
				}

				if errorBody.Error == "" {
					t.Error("error response must contain non-empty error field")
				}

				return
			}

			var body shortenResponse

			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf("cannot decode response: %v", err)
			}

			if body.OriginalURL != "https://example.com/very/long/path" {
				t.Errorf("original_url = %q, want expected URL", body.OriginalURL)
			}

			if len(body.ShortURL) != 8 {
				t.Errorf("short_url length = %d, want 8", len(body.ShortURL))
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	shortener := NewURLShortener()
	originalURL := "https://example.com/very/long/path?query=value"

	shortID, err := shortener.Shorten(originalURL)
	if err != nil {
		t.Fatalf("Shorten() unexpected error: %v", err)
	}

	tests := []struct {
		name         string
		method       string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "redirect for existing short ID",
			method:       http.MethodGet,
			path:         "/" + shortID,
			wantStatus:   http.StatusFound,
			wantLocation: originalURL,
		},
		{
			name:       "unknown short ID",
			method:     http.MethodGet,
			path:       "/unknown1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "empty short ID",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "short ID with slash",
			method:     http.MethodGet,
			path:       "/one/two",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "unsupported method",
			method:     http.MethodPost,
			path:       "/" + shortID,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	handler := redirectHandler(shortener)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != tt.wantStatus {
				t.Fatalf(
					"status = %d, want %d; body = %s",
					response.StatusCode,
					tt.wantStatus,
					recorder.Body.String(),
				)
			}

			if tt.wantLocation != "" {
				location := response.Header.Get("Location")

				if location != tt.wantLocation {
					t.Errorf(
						"Location = %q, want %q",
						location,
						tt.wantLocation,
					)
				}
			}
		})
	}
}
