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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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

	mockDB := new(dbrepo.MockDBRepo)

	expectedApp := models.ThisApp{
		ID: 1, NewApp: models.NewApp{Name: "App 1"},
	}
	app := &AuthServerApp{
		DB: mockDB,
	}

	mockDB.On("ThisApp", 1).Return(&expectedApp, nil)

	req := httptest.NewRequest(http.MethodGet, "/app/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.ThisApp
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedApp, response)
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

	expected := `{"error":true,"message":"strconv.Atoi: parsing \"\": invalid syntax"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestGetAppHandler_NotFound tests the GetApp handler when the app is not found in the database
func TestGetAppHandler_NotFound(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)
	mockDB.On("ThisApp", 1).Return(&models.ThisApp{}, errors.New("app not found"))

	app := &AuthServerApp{
		DB: mockDB,
	}

	req, err := http.NewRequest(http.MethodGet, "/app/notfound", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":true,"message":"app not found"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestThisAppHandler tests the ThisApp handler by sending a GET request to the /this-app/{id} endpoint
func TestThisAppHandler(t *testing.T) {

	mockDB := new(dbrepo.MockDBRepo)
	expectedApp := models.ThisApp{ID: 1, NewApp: models.NewApp{Name: "App 1"}}

	mockDB.On("ThisApp", 1).Return(&expectedApp, nil)
	app := &AuthServerApp{
		DB: mockDB,
	}

	req := httptest.NewRequest(http.MethodGet, "/this-app/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ThisApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.ThisApp
	err := json.NewDecoder(rr.Body).Decode(&response)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedApp, response)
}

// TestThisAppForEditHandler tests the ThisAppForEdit handler by sending a GET request to the /this-app-for-edit/{id} endpoint
// and checking the response for the expected app details.
// This handler is typically used for admin users to fetch app details for editing purposes.
func TestThisAppForEditHandler(t *testing.T) {
	mockDB := new(dbrepo.MockDBRepo)
	expectedApp := models.ThisApp{ID: 1, NewApp: models.NewApp{Name: "App 1"}}

	mockDB.On("ThisApp", 1).Return(&expectedApp, nil)
	app := &AuthServerApp{
		DB: mockDB,
	}

	req, err := http.NewRequest(http.MethodGet, "/this-app-for-edit/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ThisAppForEdit)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.ThisApp

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedApp, response)
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
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
	}

	thisApp := models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    newApp.Name,
			Release: newApp.Release,
			Path:    newApp.Path,
			Init:    newApp.Init,
			Web:     newApp.Web,
			Title:   newApp.Title,
		},
	}

	mockDB.On("InsertApp", newApp).Return(1, nil)
	app := &AuthServerApp{
		DB: mockDB,
	}

	payloadBytes, err := json.Marshal(thisApp)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/insert-app", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "AuthServerApp/json")
	req = mux.SetURLVars(req, map[string]string{"name": newApp.Name})

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
			Created: time.Now().Unix(),
			Updated: time.Now().Unix(),
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
			Created: time.Now().Unix(),
			Updated: time.Now().Unix(),
		},
	}

	mockDB.On("ThisApp", existingApp.ID).Return(&existingApp, nil)
	mockDB.On("UpdateApp", updatedApp).Return(nil)
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
	req.Header.Set("Content-Type", "AuthServerApp/json")
	req = mux.SetURLVars(req, map[string]string{
		"id":      fmt.Sprintf("%d", existingApp.ID),
		"name":    existingApp.Name,
		"release": existingApp.Release,
		"path":    existingApp.Path,
		"init":    existingApp.Init,
		"web":     existingApp.Web,
		"title":   existingApp.Title,
		"created": fmt.Sprintf("%d", existingApp.Created),
		"updated": fmt.Sprintf("%d", existingApp.Updated),
	})

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
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.DeleteApp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	expected := `{"error":false,"message":"app deleted"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
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
