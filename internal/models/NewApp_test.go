package models_test

import (
	"authserver-backend/internal/models"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewApp tests the NewApp struct for JSON serialization and deserialization.
func TestNewApp(t *testing.T) {
	// Create an instance of NewApp
	app := models.NewApp{
		Name:    "TestApp",
		Release: "1.0.0",
		Path:    "/test/path",
		Init:    "init.sh",
		Web:     "http://testapp.com",
		Title:   "Test AuthServerApp",
		Created: 1660000000, // Example timestamp
		Updated: 1660000001,
	}

	// Serialize to JSON when sending data to an external source, i.e., an API response.
	// This will ensure that the struct can be correctly converted to JSON format.
	data, err := json.Marshal(app)
	assert.NoError(t, err)

	expectedJSON := `{"name":"TestApp","release":"1.0.0","path":"/test/path","init":"init.sh","web":"http://testapp.com","title":"Test AuthServerApp","created":1660000000,"updated":1660000001}`
	assert.JSONEq(t, expectedJSON, string(data))

	// Test deserialization (unmarshal)
	// This will ensure that the struct can be correctly reconstructed from JSON
	// when coming from an external source, i.e., an API request.

	var newAppFromJSON models.NewApp
	err = json.Unmarshal(data, &newAppFromJSON)
	assert.NoError(t, err)
	assert.Equal(t, app, newAppFromJSON) // Assert equality between the original and unmarshalled structs
}

func TestNewApp_Error(t *testing.T) {
	app := models.NewApp{}

	// Calling the Error method should panic
	assert.Panics(t, func() {
		panic(app.Error())
	})
}
