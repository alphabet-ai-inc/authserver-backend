package models_test

import (
	"backend/internal/models"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThisApp(t *testing.T) {
	// Create an instance of ThisApp
	app := models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "Test App",
			Release: "1.0.0",
			Path:    "/test/path",
			Init:    "init.sh",
			Web:     "http://testapp.com",
			Title:   "Test Autserverapp",
			Created: 1660000000, // Example timestamp
			Updated: 1660000001,
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(app)
	assert.NoError(t, err)

	expectedJSON := `{"id":1,"name":"Test App","release":"1.0.0","path":"/test/path","init":"init.sh","web":"http://testapp.com","title":"Test Autserverapp","created":1660000000,"updated":1660000001}`
	assert.JSONEq(t, expectedJSON, string(data))

	// Test deserialization (unmarshal)
	var appFromJSON models.ThisApp
	err = json.Unmarshal(data, &appFromJSON)
	assert.NoError(t, err)
	assert.Equal(t, app, appFromJSON) // Assert equality between the original and unmarshalled structs
}

func TestThisApp_Error(t *testing.T) {
	app := models.ThisApp{}

	// Calling the Error method should panic
	assert.Panics(t, func() {
		panic(app.Error())
	})
}
