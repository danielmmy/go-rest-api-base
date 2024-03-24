package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gorest/internal/tools"
)

// writeResponse is a helper function to write an http response.
// it returns any possible errors.
func (app *App) writeResponse(w http.ResponseWriter, code int, data any, headers ...http.Header) error {
	var payload []byte
	var err error

	// if data is not empty prepare payload
	if data != nil {
		// if response is an error convert to marshable type with msg field
		if e, ok := data.(error); ok {
			data = &struct {
				Message string `json:"msg"`
			}{e.Error()}
		}

		payload, err = json.Marshal(data)
		if err != nil {
			tools.ErrorLogger.Println(err)
			return err
		}
	}

	// prepare response headers
	w.Header().Set("Content-Type", "application/json")
	for _, header := range headers {
		for k, v := range header {
			w.Header()[k] = v
		}
	}
	w.WriteHeader(code)

	// respond
	if _, err := w.Write(payload); err != nil {
		tools.ErrorLogger.Println(err)
		return err
	}

	return nil
}

// readJson is a helper function to read json payloads into data.
// it returns any possible errors.
func (app *App) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	// limits the size read to 1MB.
	var maxBytes int64 = 1048576 // 1MB.
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	dec := json.NewDecoder(r.Body)

	// read json into data.
	if err := dec.Decode(data); err != nil {
		tools.ErrorLogger.Println(err)
		return err
	}

	// check if body is empty.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("only one json body allowed")
	}

	return nil
}
