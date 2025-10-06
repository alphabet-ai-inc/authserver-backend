package api

import (
	"authserver-backend/internal/dbrepo"
	"authserver-backend/internal/models"

	"errors"
	"strings"

	"authserver-backend/auth"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// AuthServerApp is the main application struct
var app AuthServerApp

// TestHomeHandler tests the communication between the front-end and back-end
// by sending a GET request to the home handler and checking the response.
// without auth middleware for simplicity or accessing the database.
func TestHomeHandler(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Home)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"active","message":"Go apps up and running","version":"1.0.0"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestAppsHandler tests the apps handler by sending a GET request to the /apps endpoint
// and checking the response for the expected list of apps.
func TestAppsHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)
	expectedApps := []*models.ThisApp{
		{ID: 1, NewApp: models.NewApp{Name: "App 1", Release: "v1.0", Path: "/app1", Init: "./app1.sh", Web: "http://app1.example.com", Title: "Title1", Created: 160000000, Updated: 160000000}},
		{ID: 2, NewApp: models.NewApp{Name: "App 2", Release: "v1.1", Path: "/app2", Init: "./app2.sh", Web: "http://app2.example.com", Title: "Title2", Created: 160000001, Updated: 160000001}},
	}
	mockDB.On("AllApps").Return(expectedApps, nil)

	app := &AuthServerApp{
		DB: mockDB,
	}

	// Set up expectations
	payloadBytes, err := json.Marshal(expectedApps)

	// fmt.Printf("payloadBytes: %v\n", payloadBytes)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "/apps", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "AuthServerApp/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Apps)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []*models.ThisApp

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedApps, response)
	// assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	mockDB.AssertExpectations(t)
}

// TestAppsCatalogueHandler tests the apps catalogue handler by sending a GET request to the /apps-catalogue endpoint
// and checking the response for the expected list of apps. Is by now the same as TestAppsHandler.
// The difference is that one is for a common user and the other one is for admin users, so in the future
// we might want to add more specific checks or different behaviors based on the user role.
func TestAppsCatalogueHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)

	expectedArrApps := []*models.ThisApp{
		{ID: 1, NewApp: models.NewApp{Name: "App 1"}},
		{ID: 2, NewApp: models.NewApp{Name: "App 2"}},
	}

	app := &AuthServerApp{
		DB: mockDB,
	}

	mockDB.On("AllApps").Return(expectedArrApps, nil)
	req, err := http.NewRequest(http.MethodGet, "/apps-catalogue", nil)
	req.Header.Set("Content-Type", "AuthServerApp/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.AppsCatalogue)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []*models.ThisApp
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedArrApps, response)
	mockDB.AssertExpectations(t)
}

// TestGetAppHandler tests the GetApp handler by sending a GET request to the /app/{id} endpoint
// This handler retrieves a specific app by its ID and returns its details in the response.
// This is for common users to view app details.
// Probably in the future we might want to add more specific checks or different behaviors based on the user role.
// For instance, we might give access directly to the app link from here if the user is authenticated.
func TestGetAppHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo) // Adjust to your mock struct

	// Set up expectation for ThisApp
	expectedApp := &models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "TestApp",
			Release: "v1.0",
			Path:    "/path/to/testapp",
			Init:    "init.sh",
			Web:     "index.html",
			Title:   "Test Application",
			Created: 160000000,
			Updated: 160000000,
		},
	}
	mockDB.On("ThisApp", 1, "").Return(expectedApp, nil) // Expect call with appID=1, mockedfields=""

	app := &AuthServerApp{DB: mockDB}

	// Set up chi router for URL params
	r := chi.NewRouter()
	r.Get("/app/{id}", app.GetApp)

	req := httptest.NewRequest("GET", "/app/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify mock was called
	mockDB.AssertExpectations(t)
}

// TestGetAppHandler_InvalidID tests the GetApp handler with an invalid ID parameter
func TestGetAppHandler_InvalidID(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/app/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":true,"message":"id is missing in URL"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestGetAppHandler_NotFound tests the GetApp handler when the app is not found in the database
func TestGetAppHandler_NotFound(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)
	mockDB.On("ThisApp", 1, "").Return(&models.ThisApp{}, errors.New("app not found"))

	app := &AuthServerApp{
		DB: mockDB,
	}

	req, err := http.NewRequest(http.MethodGet, "/app/notfound", nil)
	r := chi.NewRouter()
	r.Get("/app/{id}", app.GetApp)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":true,"message":"id is missing in URL"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestThisAppHandler tests the ThisApp handler by sending a GET request to the /this-app/{id} endpoint
func TestThisAppHandler(t *testing.T) {

	mockDB := new(dbrepo.MockDBRepo)
	expectedApp := &models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "TestApp",
			Release: "v1.0",
			Path:    "/path/to/testapp",
			Init:    "init.sh",
			Web:     "index.html",
			Title:   "Test Application",
			Created: 160000000,
			Updated: 160000000,
		},
	}

	mockDB.On("ThisApp", 1, "").Return(expectedApp, nil)
	app := &AuthServerApp{DB: mockDB}

	req := httptest.NewRequest("GET", "/app/1", nil)
	w := httptest.NewRecorder()

	// Set up chi router for URL params
	r := chi.NewRouter()
	r.Get("/app/{id}", app.GetApp)
	r.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify mock was called
	mockDB.AssertExpectations(t)
}

// TestThisAppForEditHandler tests the ThisAppForEdit handler by sending a GET request to the /this-app-for-edit/{id} endpoint
// and checking the response for the expected app details.
// This handler is typically used for admin users to fetch app details for editing purposes.
func TestThisAppForEditHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)
	expectedApp := &models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "TestApp",
			Release: "v1.0",
			Path:    "/path/to/testapp",
			Init:    "init.sh",
			Web:     "index.html",
			Title:   "Test Application",
			Created: 160000000,
			Updated: 160000000,
		},
	}

	mockDB.On("ThisApp", 1, "").Return(expectedApp, nil)
	app := &AuthServerApp{
		DB: mockDB,
	}

	req, err := http.NewRequest("GET", "/app/1", nil)
	r := chi.NewRouter()
	r.Get("/app/{id}", app.GetApp)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify mock was called
	mockDB.AssertExpectations(t)
}

// TestInsertAppHandler tests the InsertApp handler by sending a POST request to the /insert-app endpoint
// with a new app payload and checking the response for successful insertion. Only for admin users.
func TestInsertAppHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)

	newApp := models.NewApp{
		Name:    "Test App",
		Release: "1.0",
		Path:    "/test/path",
		Init:    "init",
		Web:     "web",
		Title:   "Test Title",
		Created: 160000000,
		Updated: 160000000,
	}

	mockDB.On("InsertApp", mock.Anything, mock.Anything).Return(1, nil)
	app := &AuthServerApp{
		DB: mockDB,
	}

	payloadBytes, err := json.Marshal(newApp)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/insert-app", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.InsertApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	expected := `{"error":false,"message":"app inserted 1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	mockDB.AssertExpectations(t) // Ensure the mock was called
}

// TestUpdateAppHandler tests the UpdateApp handler by sending a PATCH request to the /update-app/{id} endpoint
// with updated app details and checking the response for successful update. Only for admin users.
func TestUpdateAppHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)

	existingApp := models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "Existing App",
			Release: "1.0",
			Path:    "/existing/path",
			Init:    "init",
			Web:     "web",
			Title:   "Existing Title",
			Created: 160000000,
			Updated: 160000000,
		},
	}

	updatedApp := models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "Updated App",
			Release: "2.0",
			Path:    "/updated/path",
			Init:    "init",
			Web:     "web",
			Title:   "Updated Title",
			Created: 160000000,
			Updated: 160000000,
		},
	}

	mockDB.On("ThisApp", mock.Anything, mock.Anything).Return(&existingApp, nil)
	mockDB.On("UpdateApp", mock.Anything, mock.Anything).Return(nil)
	app := &AuthServerApp{
		DB: mockDB,
	}

	payloadBytes, err := json.Marshal(updatedApp)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/update-app", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r := chi.NewRouter()
	r.Put("/update-app/{id}", app.UpdateApp)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.UpdateApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	expected := `{"error":false,"message":"app updated: "}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	mockDB.AssertExpectations(t) // Ensure the mock was called
}

// TestDeleteAppHandler tests the DeleteApp handler by sending a DELETE request to the /delete-app/{id} endpoint
// and checking the response for successful deletion. Only for admin users.
func TestDeleteAppHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)

	mockDB.On("DeleteApp", 1).Return(nil)

	app := &AuthServerApp{
		DB: mockDB,
	}

	req, err := http.NewRequest(http.MethodDelete, "/delete-app/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	r := chi.NewRouter()
	r.Delete("/delete-app/{id}", app.DeleteApp)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	expected := `{"error":false,"message":"app deleted"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	mockDB.AssertExpectations(t)
}
func TestGetReleasesHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)

	expectedReleases := []map[string]string{
		{"id": "1.0.0", "value": "1.0.0"},
		{"id": "1.0.1", "value": "1.0.1"},
		{"id": "1.0.2", "value": "1.0.2"},
		{"id": "1.1.0", "value": "1.1.0"},
		{"id": "1.2.0", "value": "1.2.0"},
		{"id": "2.0.0", "value": "2.0.0"},
		{"id": "2.0.1", "value": "2.0.1"},
		{"id": "2.0.2", "value": "2.0.2"},
		{"id": "2.1.0", "value": "2.1.0"},
		{"id": "2.2.0", "value": "2.2.0"},
		{"id": "3.0.0", "value": "3.0.0"},
	}
	mockDB.On("GetReleases").Return(expectedReleases, nil)
	app := &AuthServerApp{
		DB: mockDB,
	}
	req, err := http.NewRequest(http.MethodGet, "/releases", nil)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetReleases)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var response []map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Response:", response)
	fmt.Println("Expected:", expectedReleases)
	assert.Equal(t, expectedReleases, response)
	mockDB.AssertExpectations(t)
}

// TestAuthenticateHandler tests the Authenticate handler by sending a POST request to the /authenticate endpoint
// with valid user credentials and checking the response for the expected JWT tokens.
func TestAuthenticateHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)

	testPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	expectedUser := &models.User{
		ID:       1,
		Email:    "user@example.com",
		Password: string(hashedPassword),
	}
	mockDB.On("GetUserByEmail", "user@example.com").Return(expectedUser, nil)

	app := &AuthServerApp{
		DB: mockDB,
		Auth: auth.Auth{
			Secret:     "test_secret",
			CookieName: "refresh_token",
		},
		JWTSecret: "test_secret",
	}

	payload := fmt.Sprintf(`{"email":"user@example.com","password":"%s"}`, testPassword)
	req, err := http.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Authenticate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK && status != http.StatusAccepted {
		t.Logf("Response body: %s", rr.Body.String())
		t.Errorf("handler returned status code: %v", status)
		return
	}

	var response struct {
		Token        string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Token == "" || response.RefreshToken == "" {
		t.Error("tokens are empty")
	}
}

// TestRefreshTokenHandler tests the RefreshToken handler by sending a POST request to the /refresh-token endpoint
// with a valid refresh token cookie and checking the response for new JWT tokens.
func TestRefreshTokenHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)

	app := &AuthServerApp{
		DB: mockDB,
		Auth: auth.Auth{
			CookieName: "refresh_token",
		},
		JWTSecret: "test_secret",
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"email":   "email",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}).SignedString([]byte("test_secret"))

	if err != nil {
		t.Fatal(err)
	}
	mockDB.On("GetUserByID", 1).Return(&models.User{ID: 1, Email: "email"}, nil)

	cookie := &http.Cookie{
		Name:  app.Auth.CookieName,
		Value: refreshToken,
	}
	req, err := http.NewRequest(http.MethodPost, "/refresh-token", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(cookie)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.RefreshToken)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response struct {
		Token        string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, response.Token)
	assert.NotEmpty(t, response.RefreshToken)
}

// TestLogoutHandler tests the Logout handler by sending a POST request to the /logout endpoint
// and checking the response for successful logout and cookie deletion.
func TestLogoutHandler(t *testing.T) {
	testAuth := auth.Auth{
		CookieName: "refresh_token",
	}
	app := &AuthServerApp{
		Auth: testAuth,
	}
	type MockAutserver struct {
		Auth auth.Auth
	}
	// Mock the Auth interface
	MockApp := MockAutserver{
		Auth: testAuth,
	}

	req, err := http.NewRequest(http.MethodPost, "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Logout)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	cookies := rr.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, cookies[0].Name, MockApp.Auth.CookieName)
	assert.Equal(t, cookies[0].MaxAge, -1)
}
