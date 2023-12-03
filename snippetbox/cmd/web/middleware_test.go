package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/johnmerga/Mastering_Go/snippetbox/internal/assert"
)

func TestSecureHeader(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	// we Initialize a new dummy http.Request because the secureHeaders middleware expects a pointer to an http.Request as its parameter.
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a mock HTTP handler that we can pass to our secureHeaders
	// middleware, which writes a 200 status code and an "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// ServeHTTP() method to execute the middleware and capture the response.
	secureHeaders(next).ServeHTTP(responseRecorder, req)
	result := responseRecorder.Result()
	tests := []struct {
		name     string
		expValue string
	}{
		{
			name:     "Content-Security-Policy",
			expValue: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			name:     "Referrer-Policy",
			expValue: "origin-when-cross-origin",
		},
		{
			name:     "X-Frame-Options",
			expValue: "deny",
		},
		{
			name:     "X-XSS-Protection",
			expValue: "0",
		},
	}

	for _, header := range tests {
		t.Run(header.name, func(t *testing.T) {
			assert.Equals(t, result.Header.Get(header.name), header.expValue)
		})
	}
	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	assert.Equals(t, result.StatusCode, http.StatusOK)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equals(t, string(body), "OK")
}
