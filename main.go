package main

import (
	"authserver-backend/api"
	"authserver-backend/auth"
	"authserver-backend/internal/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx  driver for database/sql
	"github.com/joho/godotenv"
)

var port int

func main() {
	// Locate the current directory
	app := api.Autserverapp{}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Construct the path to the .env file
	envPath := filepath.Join(dir, ".env")

	// Check if the .env file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Println(".env file not found, using default values")
	} else {
		log.Println(".env file found, loading environment variables")
		// Load environment variables from .env file if present
		err = godotenv.Load(envPath)
		if err != nil {
			log.Fatal("Error loading .env file. It is not readable.")
		}
	}
	app.DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_EXTERNAL_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	app.JWTSecret = os.Getenv("JWT_SECRET")

	// Set default values for the app configuration
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.IntVar(&port, "port", 8080, "API server port")

	flag.Parse()
	// Initialize the database connection
	if app.DSN == "" {
		log.Fatal("DSN environment variable is not set")
	}
	if app.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	log.Println("DSN: " + app.DSN)

	repo := &dbrepo.PostgresDBRepo{}
	db, err := repo.ConnectToDB(app.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// Test the database connection

	if err = db.Ping(); err != nil {
		log.Fatalf("Database not connected: %v", err)
	}

	log.Println("Connected to the database")

	// Set up the database connection pool
	defer db.Close()

	// Set the DB field in repo
	repo.DB = db

	// Assign the initialized repo to app.DB
	app.DB = repo

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

	// Set up your app's router/handler
	handler := app.Routes()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
