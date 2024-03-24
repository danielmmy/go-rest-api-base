package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"gorest/api"
	"gorest/internal/tools"

	"github.com/go-chi/chi/v5"
)

var readJsonAlias = (*App).readJson
var writeResponseAlias = (*App).writeResponse
var repository = tools.NewFederationRepository

func (app *App) addFederation(w http.ResponseWriter, r *http.Request) {
	// read federation data from body
	federation := new(api.Federation)
	if err := readJsonAlias(app, w, r, federation); err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusBadRequest, err)
		return
	}

	// connect to repo
	repo, err := repository()
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	// insert federation and respond
	code, err := (*repo).AddFederation(federation)
	if err := writeResponseAlias(app, w, code, err); err != nil {
		tools.ErrorLogger.Println(err)
	}
}

func (app *App) GetFederation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusBadRequest, err)
		return
	}

	repo, err := repository()
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	fed := (*repo).GetFederation(id)
	if fed == nil {
		if err := writeResponseAlias(app, w, http.StatusNotFound, fmt.Errorf("federation %d not found", id)); err != nil {
			tools.ErrorLogger.Println(err)
		}
		return
	}

	// example on how to change response headers
	header := http.Header{
		"Content-Type": {"text/plain"},
	}
	if err := writeResponseAlias(app, w, http.StatusOK, fed, header); err != nil {
		tools.ErrorLogger.Println(err)
	}
}

func (app *App) getFederations(w http.ResponseWriter, r *http.Request) {
	repo, err := repository()
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	federations := (*repo).GetFederations()
	if err := writeResponseAlias(app, w, http.StatusOK, federations); err != nil {
		tools.ErrorLogger.Println(err)
	}
}

func (app *App) updateFederation(w http.ResponseWriter, r *http.Request) {
	fed := new(api.Federation)
	if err := readJsonAlias(app, w, r, fed); err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusBadRequest, err)
		return
	}

	fed.Id = id

	repo, err := repository()
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	code, err := (*repo).UpdateFederation(fed)
	if err := writeResponseAlias(app, w, code, err); err != nil {
		tools.ErrorLogger.Println(err)
	}
}

func (app *App) deleteFederation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusBadRequest, err)
		return
	}

	repo, err := repository()
	if err != nil {
		tools.ErrorLogger.Println(err)
		writeResponseAlias(app, w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	code, err := (*repo).DeleteFederation(id)
	if err := writeResponseAlias(app, w, code, err); err != nil {
		tools.ErrorLogger.Println(err)
	}
}
