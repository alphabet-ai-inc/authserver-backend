// Package api provides HTTP routing and middleware for the authserver-backend application.
package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// EnableCORS is a middleware function that enables Cross-Origin Resource Sharing (CORS).
// It allows requests from specified origins and handles preflight OPTIONS requests.
// Adjust the allowed origins as needed for your application.
// For production, ensure to set the ALLOWED_ORIGINS environment variable appropriately.

func (app *AuthServerApp) EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		origin := r.Header.Get("Origin")
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

		found := false

		for _, allowedOrigin := range strings.Split(allowedOrigins, ",") {
			if origin == strings.TrimSpace(allowedOrigin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				found = true
				break
			}
		}
		if !found && origin != "" {
			log.Printf("Origin not allowed: %s (Allowed: %s)", origin, allowedOrigins)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// for production, needed to change localhost here and move the following line
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		// deleting it from OPTIONS.

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
func (app *AuthServerApp) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// inspecting the front-end request
		// log.Printf("Request: %s %s", r.Method, r.URL.Path)
		// log.Printf("Headers: %v", r.Header)
		// log.Printf("Body: %v", r.Body)

		_, _, err := app.Auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("El usuario no está autorizado: %d", http.StatusUnauthorized)))
			return
		}
		next.ServeHTTP(w, r)
	})
}
