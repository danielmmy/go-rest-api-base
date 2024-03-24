package handlers

import (
	"net/http"

	"gorest/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *App) NewHandler() http.Handler {
	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(chimiddle.StripSlashes)
	mux.Use(chimiddle.Heartbeat("/health-check"))
	mux.Route("/federations", func(r chi.Router) {
		r.Use(middleware.Authorize)
		r.Post("/", app.addFederation)
		r.Get("/{id}", app.GetFederation)
		r.Get("/", app.getFederations)
		r.Put("/{id}", app.updateFederation)
		r.Delete("/{id}", app.deleteFederation)
	})

	return mux
}
