// Package api provides HTTP routing and middleware for the authserver-backend application.
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Routes sets up the application's HTTP routes and middleware.
// It returns an http.Handler that can be used by the HTTP server.
//
// The following endpoints are registered:
//   - GET    /                  : Home
//   - POST   /authenticate      : Authenticate user and issue JWT
//   - GET    /refresh           : Refresh JWT token
//   - GET    /logout            : Log out user
//   - POST   /validatesession   : Validate JWT session
//   - GET    /apps              : List apps
//   - GET    /apps/{id}         : Get app by ID
//
// The /admin subrouter is protected by authentication middleware and provides:
//   - GET    /admin/apps              : List all apps (admin)
//   - GET    /admin/apps/{id}         : Get app for editing (admin)
//   - POST   /admin/apps/0            : Insert new app (admin)
//   - PATCH  /admin/apps/{id}         : Update app (admin)
//   - DELETE /admin/apps/{id}         : Delete app (admin)

func (app *AuthServerApp) Routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.EnableCORS)

	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.Authenticate)
	mux.Get("/refresh", app.RefreshToken)
	mux.Get("/logout", app.Logout)
	mux.Post("/validatesession", app.ValidateSession)
	mux.Get("/apps", app.Apps)
	mux.Get("/apps/{id}", app.GetApp)
	mux.Get("/releases", app.GetReleases) // Assuming app is your AuthServerApp instance
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/apps", app.AppsCatalogue)
		mux.Get("/apps/{id}", app.ThisAppForEdit)
		mux.Post("/apps/0", app.InsertApp)
		mux.Patch("/apps/{id}", app.UpdateApp)
		mux.Delete("/apps/{id}", app.DeleteApp)

	})

	return mux
}
