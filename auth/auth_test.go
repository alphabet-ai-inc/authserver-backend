package auth_test

import (
	"authserver-backend/auth"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGenerateTokenPair tests the GenerateTokenPair method of the Auth struct.
func TestGenerateTokenPair(t *testing.T) {

	authService := auth.Auth{
		Issuer:        "testIssuer",
		Audience:      "testAudience",
		Secret:        "testSecret",
		TokenExpiry:   time.Minute, // 1 minute token expiry
		RefreshExpiry: time.Hour,   // 1 hour refresh token expiry
	}

	var user = auth.JWTUser{
		ID:    1,
		Email: "admin@example.com",
	}

	tokenPairs, err := authService.GenerateTokenPair(&user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tokenPairs.Token == "" || tokenPairs.RefreshToken == "" {
		t.Fatalf("Expected non-empty tokens, got %+v", tokenPairs)
	}
}

// Instead of using a static token, we generate one for testing
func TestGetTokenFromHeaderAndVerify(t *testing.T) {

	authService := auth.Auth{
		Issuer:        "testIssuer",
		Audience:      "testAudience",
		Secret:        "testSecret",
		TokenExpiry:   time.Minute, // 1 minute token expiry
		RefreshExpiry: time.Hour,   // 1 hour refresh token expiry
	}

	var user = auth.JWTUser{
		ID:    1,
		Email: "admin@example.com",
	}

	// Generate a token
	tokenPairs, err := authService.GenerateTokenPair(&user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	t.Logf("Generated token pairs: %+v", tokenPairs)
	// Create a request with the authorization header
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPairs.Token)

	rr := httptest.NewRecorder() // Response recorder to capture the response

	// Call the method to verify token
	token, claims, err := authService.GetTokenFromHeaderAndVerify(rr, req)
	t.Logf("token: %+v", token)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("Expected non-empty token, got %s", token)
	}

	if claims.Issuer != authService.Issuer {
		t.Fatalf("Expected issuer to be %s, got %s", authService.Issuer, claims.Issuer)
	}
}

// TestGetRefreshCookie tests the GetRefreshCookie method of the Auth struct.
func TestGetRefreshCookie(t *testing.T) {
	auth := &auth.Auth{
		CookieName:    "refresh_token",
		CookiePath:    "/",
		CookieDomain:  "localhost",
		RefreshExpiry: time.Hour,
	}

	refreshToken := "exampleRefreshToken"
	cookie := auth.GetRefreshCookie(refreshToken)

	if cookie.Name != auth.CookieName {
		t.Errorf("Expected cookie name %s, got %s", auth.CookieName, cookie.Name)
	}

	if cookie.Value != refreshToken {
		t.Errorf("Expected cookie value %s, got %s", refreshToken, cookie.Value)
	}
}

// TestGetExpiredRefreshCookie tests the GetExpiredRefreshCookie method of the Auth struct.
func TestGetExpiredRefreshCookie(t *testing.T) {
	auth := &auth.Auth{
		CookieName:   "refresh_token",
		CookiePath:   "/",
		CookieDomain: "localhost",
	}

	cookie := auth.GetExpiredRefreshCookie()

	if cookie.Value != "" {
		t.Errorf("Expected cookie value to be empty, got %s", cookie.Value)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("Expected MaxAge to be -1 for expired cookie, got %d", cookie.MaxAge)
	}
}

// TestGetTokenFromHeaderInvalidFormat tests the GetTokenFromHeaderAndVerify method with an invalid token format.
func TestGetTokenFromHeaderInvalidFormat(t *testing.T) {
	authService := &auth.Auth{
		Issuer:        "testIssuer",
		Audience:      "testAudience",
		Secret:        "testSecret",
		JWTSecret:     "testSecret",
		CookieName:    "refresh_token",
		CookiePath:    "/",
		CookieDomain:  "localhost",
		RefreshExpiry: time.Hour,
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidTokenFormat")

	rr := httptest.NewRecorder()

	// Call the method to verify token
	_, _, err := authService.GetTokenFromHeaderAndVerify(rr, req)

	if err == nil || err.Error() != "invalid auth header" {
		t.Fatalf("Expected invalid auth header error, got %v", err)
	}
}

// TestGetTokenFromHeaderNoAuth tests the GetTokenFromHeaderAndVerify method with no Authorization header.
func TestGetTokenFromHeaderNoAuth(t *testing.T) {
	authService := &auth.Auth{
		JWTSecret:     "testSecret",
		CookieName:    "refresh_token",
		CookiePath:    "/",
		CookieDomain:  "localhost",
		RefreshExpiry: time.Hour,
	}

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Call the method to verify token
	_, _, err := authService.GetTokenFromHeaderAndVerify(rr, req)

	if err == nil || err.Error() != "no auth header" {
		t.Fatalf("Expected no auth header error, got %v", err)
	}
}
