package main

import (
	"backend/api"
	"backend/internal/dbrepo"
	"backend/pkg/auth"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

func main() {
	// set application config
	var app api.Application

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=jpassano password=jP1732 dbname=autserver sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")

	flag.Parse()
	// Initialize the database connection
	repo := &dbrepo.PostgresDBRepo{}
	app.DB = repo

	var err error
	db, err := app.DB.ConnectToDB(app.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Connected to the database")

	// Set up the database connection pool
	defer db.Close()

	// Assign the database repo to app.DB
	app.DB = &dbrepo.PostgresDBRepo{DB: db}

	app.Auth = auth.Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	// Start a web server
	fmt.Printf("Starting server on port %d\n", port)

	// Set up your application's router/handler
	handler := app.Routes()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))

}
