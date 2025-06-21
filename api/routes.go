package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Autserverapp) Routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.EnableCORS)

	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.Authenticate)
	mux.Get("/refresh", app.RefreshToken)
	mux.Get("/logout", app.Logout)

	mux.Get("/apps", app.Apps)
	mux.Get("/apps/{id}", app.GetApp)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/apps", app.AppsCatalogue)
		mux.Get("/apps/{id}", app.ThisAppForEdit)
		mux.Put("/apps/0", app.InsertApp)
		mux.Patch("/apps/{id}", app.UpdateApp)
		mux.Delete("/apps/{id}", app.DeleteApp)

	})

	return mux
}
