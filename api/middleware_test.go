package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// type mockAuth struct {
// 	// Implement the interface
// 	GetTokenFromHeaderAndVerifyFunc func(w http.ResponseWriter, r *http.Request) (string, *Claims, error)
// }

// // Implement the method for the mocked auth
// func (m *mockAuth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
// 	return m.GetTokenFromHeaderAndVerifyFunc(w, r)
// }

// type mockAutserverapp struct {
// 	auth AuthInterface // Use the Auth interface defined earlier
// }

// func (m *mockAutserverapp) authRequired(nextHandler http.HandlerFunc) any {
// 	panic("unimplemented")
// }

// Test for CORS Middleware

func TestEnableCORS(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://example.com")

	// Create a mock HTTP handler
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := app.EnableCORS(mockHandler)

	// Test a GET request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	corsHandler.ServeHTTP(recorder, req)

	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	// assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))

	// Test an OPTIONS request for preflight
	reqOptions := httptest.NewRequest(http.MethodOptions, "/", nil)
	recorderOptions := httptest.NewRecorder()

	corsHandler.ServeHTTP(recorderOptions, reqOptions)

	respOptions := recorderOptions.Result()
	assert.Equal(t, http.StatusOK, respOptions.StatusCode)
	assert.Equal(t, "true", respOptions.Header.Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", respOptions.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Accept, Content-Type, X-CSRF-Token, Authorization", respOptions.Header.Get("Access-Control-Allow-Headers"))
}

// Test for Authentication Middleware
// func TestAuthRequired(t *testing.T) {
// 	// Create the mock auth instance
// 	mockAuth := &mockAuth{
// 		GetTokenFromHeaderAndVerifyFunc: func(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
// 			return "token", &Claims{jwt.RegisteredClaims{Subject: "1"}}, nil
// 		},
// 	}

// 	app := &mockAutserverapp{
// 		auth: mockAuth, // Use the mock auth
// 	}

// 	// Create a handler to test the next middleware
// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK) // Return 200 OK for the next handler
// 	})

// 	authHandler := app.authRequired(nextHandler)

// 	// Test authorized request
// 	reqAuth := httptest.NewRequest(http.MethodGet, "/", nil)
// 	recorderAuth := httptest.NewRecorder()

// 	authHandler.ServeHTTP(recorderAuth, reqAuth)

// 	// Assert the response for the authorized request
// 	assert.Equal(t, http.StatusOK, recorderAuth.Code)

// 	// Now test an unauthorized request
// 	mockAuth.GetTokenFromHeaderAndVerifyFunc = func(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
// 		return "", nil, errors.New("unauthorized") // Simulate an error
// 	}

// 	reqUnauthorized := httptest.NewRequest(http.MethodGet, "/", nil)
// 	recorderUnauthorized := httptest.NewRecorder()

// 	authHandler.ServeHTTP(recorderUnauthorized, reqUnauthorized)

// 	// Assert the response for the unauthorized request
// 	assert.Equal(t, http.StatusUnauthorized, recorderUnauthorized.Code)
// }
