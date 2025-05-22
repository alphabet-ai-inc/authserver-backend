package api_test

import (
	"backend/api"
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"strings"

	"backend/pkg/auth"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// var mockDB = new(repository.MockDBRepo)
// var app api.Application

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app := api.Application{
		Auth:      auth.Auth{},
		JWTSecret: "secret",
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
func TestAppsHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)
	app := api.Application{
		DB: mockDB,
	}

	expectedApps := []*models.ThisApp{
		{ID: 1, Name: "App 1"},
		{ID: 2, Name: "App 2"},
	}

	mockDB.On("AllApps").Return(expectedApps, nil)

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

	req.Header.Set("Content-Type", "application/json")

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
	mockDB.AssertExpectations(t)
}

func TestAuthenticateHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)
	mockPassMatch := new(models.MockPassMatch)

	app := api.Application{
		DB:          mockDB,
		CheckPasswd: mockPassMatch,
	}

	// Set up expectations
	// mockDB.On("GetUserByEmail", "user@example.com").Return(models.User{ID: 1, Email: "user@example.com"}, nil)
	// mockDB.On("PasswordMatches", "password").Return(true, nil)

	// Mock the GetUserByEmail method
	mockDB.On("GetUserByEmail", "user@example.com").Return(&models.User{
		ID:       1,
		Email:    "user@example.com",
		Password: "hashedpassword", // Replace with a valid hashed password
	}, nil)

	mockPassMatch.On("PasswordMatches", "password").Return(true, nil)

	payload := `{"email": "user@example.com", "password": "password"}`

	req, err := http.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(payload))

	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Authenticate)

	handler.ServeHTTP(rr, req)

	rr.Code = httptest.NewRecorder().Result().StatusCode

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

func TestRefreshTokenHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)

	mockAuth := auth.Auth{
		MockToken:        "token",
		MockRefreshToken: "refresh_token",
		CookieName:       "refresh_token",
		JWTSecret:        "secret",
	}
	app := api.Application{
		DB:   mockDB,
		Auth: mockAuth,
	}

	// Create a valid refresh token
	refreshToken, err := app.Auth.GenerateRefreshToken(&auth.JWTUser{ID: 1, Email: "email"})
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

func TestLogoutHandler(t *testing.T) {
	mockAuth := auth.Auth{
		CookieName: "refresh_token",
	}
	app := api.Application{
		Auth: mockAuth,
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
	assert.Equal(t, cookies[0].Name, app.Auth.CookieName)
	assert.Equal(t, cookies[0].MaxAge, -1)
}

func TestAppsCatalogueHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)
	app := api.Application{
		DB: mockDB,
	}

	expectedArrApps := []*models.ThisApp{
		{ID: 1, Name: "App 1"},
		{ID: 2, Name: "App 2"},
	}

	mockDB.On("AllApps").Return(expectedArrApps, nil)
	req, err := http.NewRequest(http.MethodGet, "/apps-catalogue", nil)
	req.Header.Set("Content-Type", "application/json")

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

func TestGetAppHandler(t *testing.T) {

	mockDB := new(repository.MockDBRepo)
	app := api.Application{
		DB: mockDB,
	}

	expectedApp := models.ThisApp{ID: 1, Name: "App 1"}
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

func TestGetAppHandler_InvalidID(t *testing.T) {
	mockDB := new(repository.MockDBRepo)

	req, err := http.NewRequest(http.MethodGet, "/app/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	app := api.Application{
		DB: mockDB,
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

func TestGetAppHandler_NotFound(t *testing.T) {
	mockDB := new(repository.MockDBRepo)
	app := api.Application{
		DB: mockDB,
	}
	mockDB.On("ThisApp", 1).Return(&models.ThisApp{}, errors.New("app not found"))

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

func TestThisAppHandler(t *testing.T) {

	mockDB := new(repository.MockDBRepo)
	app := api.Application{
		DB: mockDB,
	}
	expectedApp := models.ThisApp{ID: 1, Name: "App 1"}

	mockDB.On("ThisApp", 1).Return(&expectedApp, nil)

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

func TestThisAppForEditHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)
	app := api.Application{
		DB: mockDB,
	}
	expectedApp := models.ThisApp{ID: 1, Name: "App 1"}

	mockDB.On("ThisApp", 1).Return(&expectedApp, nil)

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
func TestInsertAppHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)

	app := api.Application{
		DB: mockDB,
	}

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
		ID:      1,
		Name:    newApp.Name,
		Release: newApp.Release,
		Path:    newApp.Path,
		Init:    newApp.Init,
		Web:     newApp.Web,
		Title:   newApp.Title,
	}

	mockDB.On("InsertApp", newApp).Return(1, nil)

	payloadBytes, err := json.Marshal(thisApp)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/insert-app", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
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
func TestUpdateAppHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)

	app := api.Application{
		DB: mockDB,
	}
	existingApp := models.ThisApp{
		ID:      1,
		Name:    "Existing App",
		Release: "1.0",
		Path:    "/existing/path",
		Init:    "init",
		Web:     "web",
		Title:   "Existing Title",
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
	}

	updatedApp := models.ThisApp{
		ID:      1,
		Name:    "Updated App",
		Release: "2.0",
		Path:    "/updated/path",
		Init:    "init",
		Web:     "web",
		Title:   "Updated Title",
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
	}

	mockDB.On("ThisApp", existingApp.ID).Return(&existingApp, nil)
	mockDB.On("UpdateApp", updatedApp).Return(nil)
	payloadBytes, err := json.Marshal(updatedApp)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/update-app", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
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

func TestDeleteAppHandler(t *testing.T) {
	mockDB := new(repository.MockDBRepo)

	app := api.Application{
		DB: mockDB,
	}
	mockDB.On("DeleteApp", 1).Return(nil)
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
