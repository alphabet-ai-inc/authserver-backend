package main

import (
	"backend/api"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/yourusername/yourproject/api"
)

func TestHome(t *testing.T) {
	// Initialize Autserverapp with necessary data

	app := api.Autserverapp{
		JWTSecret: "testsecret",
	}

	// Create a new HTTP request for the home route
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Call the Home handler
	app.Home(w, req)

	// Check the response code
	assert.Equal(t, http.StatusOK, w.Code)

	// Optional: Check the response body if necessary
	expectedResponse := `{"status":"active","message":"Go apps up and running","version":"1.0.0"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

// If you have other routes set up, you can similarly test them.

func TestMain(m *testing.M) {
	// You could perform setup here if necessary

	// Run the tests
	m.Run()

	// Any teardown logic if needed
}
