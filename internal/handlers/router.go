package handlers

import (
	"net/http"

	"gorest/internal/middleware"
)

func (app *App) NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	federationRouter := routeGroup{
		basePath: "/federations",
		ServeMux: mux,
	}
	federationRouter.Use(middleware.Authorize)
	federationRouter.HandleFunc(http.MethodPost, "", app.addFederation)
	federationRouter.HandleFunc(http.MethodGet, "/{id}", app.GetFederation)
	federationRouter.HandleFunc(http.MethodGet, "", app.getFederations)
	federationRouter.HandleFunc(http.MethodPut, "/{id}", app.updateFederation)
	federationRouter.HandleFunc(http.MethodDelete, "/{id}", app.deleteFederation)

	return mux
}
