package utils_test

import (
	"authserver-backend/internal/utils"

	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWriteJSON tests the WriteJSON method of the JSONResponse struct.
func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	err := errors.New("test error")
	utils.JSONResponse{}.ErrorJSON(rr, err, http.StatusBadRequest)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":true,"message":"test error"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestReadJSON(t *testing.T) {

	// Create a sample request body
	requestBody := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(requestBody))
	recorder := httptest.NewRecorder()

	var payload struct {
		Name string `json:"name"`
	}

	// Call the readJSON method
	var jsr utils.JSONResponse
	err := jsr.ReadJSON(recorder, req, &payload)

	// Check for errors
	assert.NoError(t, err)
	assert.Equal(t, "test", payload.Name)
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	var jsr utils.JSONResponse

	// Create a request with invalid JSON body
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":`))
	recorder := httptest.NewRecorder()

	var payload struct {
		Name string `json:"name"`
	}

	// Call the readJSON method
	err := jsr.ReadJSON(recorder, req, &payload)

	// Check for the expected error
	assert.Error(t, err)
	assert.Equal(t, "unexpected EOF", err.Error())
}

func TestReadJSON_TooManyJSONValues(t *testing.T) {
	var jsr utils.JSONResponse

	// Create a request with multiple JSON values
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"test"}{"name":"test2"}`))
	recorder := httptest.NewRecorder()

	var payload struct {
		Name string `json:"name"`
	}

	// Call the readJSON method
	err := jsr.ReadJSON(recorder, req, &payload)

	// Check for the expected error
	assert.Error(t, err)
	assert.Equal(t, "body must only contain a single JSON value", err.Error())
}

func TestErrorJSON(t *testing.T) {
	var jsr utils.JSONResponse
	recorder := httptest.NewRecorder()

	err := errors.New("test error message")

	// Call errorJSON with an error
	errResponse := jsr.ErrorJSON(recorder, err)

	// Check for errors in errorJSON method itself
	assert.NoError(t, errResponse)

	// Check response code
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Check response body
	var jsonResponse utils.JSONResponse
	err = json.NewDecoder(recorder.Body).Decode(&jsonResponse)
	assert.NoError(t, err)
	assert.True(t, jsonResponse.Error)
	assert.Equal(t, "test error message", jsonResponse.Message)
}
