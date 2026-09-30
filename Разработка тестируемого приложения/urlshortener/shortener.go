package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sync"
)

var (
	ErrInvalidURL = errors.New("invalid URL: expected absolute HTTP or HTTPS URL")
	ErrNotFound   = errors.New("short URL not found")
)

type URLShortener struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

// Shorten validates an original URL, generates a unique short ID,
// and stores the relation in memory.
func (us *URLShortener) Shorten(originalURL string) (string, error) {
	if !isValidURL(originalURL) {
		return "", ErrInvalidURL
	}

	for {
		shortID := generateShortID()

		us.mu.Lock()

		_, alreadyExists := us.urls[shortID]
		if !alreadyExists {
			us.urls[shortID] = originalURL
			us.mu.Unlock()

			return shortID, nil
		}

		us.mu.Unlock()
	}
}

// GetOriginal returns the original URL associated with a short ID.
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	if shortID == "" {
		return "", ErrNotFound
	}

	us.mu.RLock()
	originalURL, exists := us.urls[shortID]
	us.mu.RUnlock()

	if !exists {
		return "", ErrNotFound
	}

	return originalURL, nil
}

// generateShortID returns an 8-character URL-safe identifier.
// Six random bytes become exactly eight Base64 Raw URL characters.
func generateShortID() string {
	bytes := make([]byte, 6)

	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("cannot generate short ID: %v", err))
	}

	return base64.RawURLEncoding.EncodeToString(bytes)
}

// isValidURL accepts only absolute HTTP or HTTPS URLs with a host.
func isValidURL(str string) bool {
	parsedURL, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return parsedURL.Host != ""
}
