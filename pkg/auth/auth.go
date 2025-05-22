package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/golang-jwt/jwt/v4"
)

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

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.JWTSecret))
}

func (j *Auth) GenerateTokenPair(user *JWTUser) (TokenPairs, error) {
	accessToken, err := j.GenerateRefreshToken(user)
	if err != nil {
		return TokenPairs{}, err
	}

	refreshToken, err := j.GenerateRefreshToken(user)
	if err != nil {
		return TokenPairs{}, err
	}

	return TokenPairs{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

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
	// Parece estar mal el siguiente header array. Debería pasarle una autorizacion el cliente con el request
	w.Header().Add("Vary", "Authorization")

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

	token := headerParts[1]

	// declare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing metdhod: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
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
