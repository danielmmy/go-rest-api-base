package handlers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// test handler health check
// should write 200 as response
func TestNewHandlerHealthcheck(t *testing.T) {
	// arrage
	app := NewApp()
	sut := app.NewHandler()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health-check", nil)
	want := http.StatusOK

	// act
	sut.ServeHTTP(w, r)

	// assert
	if w.Code != want {
		t.Fatalf("ServeHTTP(w, r) = %d want %d", w.Code, want)
	}
}

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

	if _, ok := handler.(*http.ServeMux); !ok {
		t.Fatalf(`NewHandler() = %q want "*chi.Mux"`, reflect.TypeOf(handler))
	}
}
