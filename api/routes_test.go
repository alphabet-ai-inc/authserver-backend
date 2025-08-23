package api_test

import (
	"authserver-backend/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRoutes verifies the routing and middleware for the Autserverapp
func TestRoutes(t *testing.T) {
	// mockDB := new(repository.MockDBRepo)

	app := api.Autserverapp{}

	// Create a test server using the routes
	ts := httptest.NewServer(app.Routes())
	defer ts.Close()

	// Test the Home route
	res, err := http.Get(ts.URL + "/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Test the Authenticate route (assuming it expects JSON)
	req, err := http.NewRequest("POST", ts.URL+"/authenticate", nil) // You may want to add a proper JSON body here.
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// Test the Refresh Token route
	req, err = http.NewRequest("GET", ts.URL+"/refresh", nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// Test the Logout route
	req, err = http.NewRequest("GET", ts.URL+"/logout", nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)

	// Test the Admin route with authentication (would require a valid token here)
	// Assuming you need to set up a mock for the authentication to test these routes.
	// Let's create a test for accessing secure routes.

	adminReq, err := http.NewRequest("GET", ts.URL+"/admin/apps", nil)
	assert.NoError(t, err)

	// Assume you have a valid JWT token you would send in the Authorization header
	adminReq.Header.Set("Authorization", "Bearer VALID_JWT_TOKEN")

	res, err = http.DefaultClient.Do(adminReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode) // or whatever your auth check returns for invalid token
}
