package handlers

import (
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
)

// test NewHandler
// should return http.Handler
func TestNewHandlerSuccess(t *testing.T) {
	// arrage
	app := NewApp()

	// act
	handler := app.NewHandler()

	// assert
	if handler == nil {
		t.Fatal("NewHandler() = <nil> want http.Handler")
	}

	if _, ok := handler.(*chi.Mux); !ok {
		t.Fatalf(`NewHandler() = %q want "*chi.Mux"`, reflect.TypeOf(handler))
	}
}
