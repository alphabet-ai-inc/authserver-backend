package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Auth struct holds configuration for JWT authentication
// It is used to generate and validate JWT tokens.
// JWT tokens are signed using the HMAC SHA256 algorithm,
// and contains methods for generating and validating tokens,
// that include user information in the token claims.
// The claims are custom and include user ID and email.
// Other claims like issuer, audience, issued at and expiry are also included.
// The struct also includes methods for generating refresh tokens
// and setting refresh tokens in HTTP cookies. The tokens expire after a configurable duration.
// The struct also includes methods for mocking tokens for testing purposes.
type Auth struct {
	Issuer           string
	Audience         string
	Secret           string
	MockToken        string
	MockRefreshToken string
	TokenExpiry      time.Duration
	RefreshExpiry    time.Duration
	CookieDomain     string
	CookiePath       string
	CookieName       string
	JWTSecret        string
}

// AuthInterface defines the methods that any authentication service should implement.
// This allows for flexibility in swapping out different authentication implementations.
// In this case, it includes the mocking of GetTokenFromHeaderAndVerify method for testing purposes.
// Then, in tests, we can create a mock struct that implements this interface and in development we can use the real Auth struct.
type AuthInterface interface {
	GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error)
}

type JWTUser struct {
	ID    int
	Email string
}

type MockAuth struct {
	Token        string
	RefreshToken string
}

// TokenPairs holds the access and refresh tokens. It is used to return both tokens together.
// In this case, both tokens are strings with the access token being the main JWT token used for authentication
// and the refresh token being used to obtain a new access token when the original expires.
type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Claims defines the custom JWT claims used in the tokens. It includes user ID and email,
// along with standard claims like issuer, issued at, and expiry time.
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func (j *Auth) GenerateRefreshToken(user *JWTUser) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			Issuer:    j.Issuer,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.JWTSecret))
}

func (j *Auth) GenerateTokenPair(user *JWTUser) (TokenPairs, error) {
	accessClaims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			Issuer:    j.Issuer,
			ExpiresAt: time.Now().Add(j.TokenExpiry).Unix(), // Используем TokenExpiry
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.JWTSecret))
	if err != nil {
		return TokenPairs{}, err
	}

	refreshTokenString, err := j.GenerateRefreshToken(user)
	if err != nil {
		return TokenPairs{}, err
	}

	return TokenPairs{
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// GetRefreshCookie creates an HTTP cookie to store the refresh token securely.
func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     "/",
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

// GetExpiredRefreshCookie creates an HTTP cookie that effectively deletes the refresh token by setting its expiration date in the past.
func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// get auth header
	authHeader := r.Header.Get("Authorization")

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}
	//split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}
	// get the token part from the header that comes after Bearer in the header, part of the request
	// that comes as parameter to this function
	token := headerParts[1]

	// declare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.JWTSecret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}
	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("invalid issuer")
	}
	return token, claims, nil
}

// GenerateRefreshToken generates a refresh token for use in testing
func (j *Auth) MockGenerateRefreshToken(user *JWTUser, secret string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["email"] = user.Email
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	refreshTokenClaims["exp"] = time.Now().UTC().Add(24 * time.Hour).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedRefreshToken, nil
}
