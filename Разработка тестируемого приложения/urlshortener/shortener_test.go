package main

import (
	"errors"
	"strings"
	"testing"
)

func TestURLShortener_Shorten(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTP URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL with path and query",
			url:     "https://google.com/search?q=test",
			wantErr: false,
		},
		{
			name:    "invalid plain text",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "empty string",
			url:     "",
			wantErr: true,
		},
		{
			name:    "unsupported FTP scheme",
			url:     "ftp://example.com/file.txt",
			wantErr: true,
		},
		{
			name:    "URL without host",
			url:     "https:///path",
			wantErr: true,
		},
		{
			name:    "relative path",
			url:     "/very/long/path",
			wantErr: true,
		},
	}

	shortener := NewURLShortener()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := shortener.Shorten(tt.url)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Shorten(%q) error = %v, wantErr = %v", tt.url, err, tt.wantErr)
			}

			if tt.wantErr {
				if shortID != "" {
					t.Errorf("Shorten(%q) shortID = %q, want empty value", tt.url, shortID)
				}

				return
			}

			if len(shortID) < 6 || len(shortID) > 8 {
				t.Errorf("short ID length = %d, want from 6 to 8", len(shortID))
			}

			if strings.ContainsAny(shortID, "+/=") {
				t.Errorf("short ID %q is not URL-safe", shortID)
			}
		})
	}
}

func TestURLShortener_GetOriginal(t *testing.T) {
	shortener := NewURLShortener()

	shortID, err := shortener.Shorten("https://example.com/very/long/path")
	if err != nil {
		t.Fatalf("Shorten() unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		shortID string
		wantURL string
		wantErr error
	}{
		{
			name:    "existing short ID",
			shortID: shortID,
			wantURL: "https://example.com/very/long/path",
			wantErr: nil,
		},
		{
			name:    "unknown short ID",
			shortID: "unknown1",
			wantURL: "",
			wantErr: ErrNotFound,
		},
		{
			name:    "empty short ID",
			shortID: "",
			wantURL: "",
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := shortener.GetOriginal(tt.shortID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"GetOriginal(%q) error = %v, want %v",
					tt.shortID,
					err,
					tt.wantErr,
				)
			}

			if gotURL != tt.wantURL {
				t.Errorf(
					"GetOriginal(%q) = %q, want %q",
					tt.shortID,
					gotURL,
					tt.wantURL,
				)
			}
		})
	}
}

func TestURLShortener_ShortenGeneratesUniqueIDs(t *testing.T) {
	shortener := NewURLShortener()
	seen := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		shortID, err := shortener.Shorten("https://example.com/path/" + string(rune('a'+i%26)))
		if err != nil {
			t.Fatalf("Shorten() unexpected error: %v", err)
		}

		if seen[shortID] {
			t.Fatalf("duplicate short ID generated: %q", shortID)
		}

		seen[shortID] = true
	}
}
