package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// test Use success
// should add middleware to pipeliine
func TestUseSuccess(t *testing.T) {
	// arrange
	sut := routeGroup{}
	want := 1

	// act
	sut.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	// assert
	if len(sut.middlewares) != want {
		t.Fatalf("Use = %d want %d", len(sut.middlewares), want)
	}
}

// test Handle Success
// should add handler to mux
func TestHandleSucess(t *testing.T) {
	// arrange
	sut := routeGroup{
		basePath: "/tests",
		ServeMux: http.NewServeMux(),
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/tests/test1", nil)
	called := false

	// act
	sut.Handle(http.MethodGet, "/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	sut.ServeHTTP(w, r)

	// assert
	if !called {
		t.Fatalf(`Handle(http.MethodGet, "/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true})) = %v want true`, called)
	}
}

// test HandleFunc Success
// should add handler func to mux
func TestHandleFuncSucess(t *testing.T) {
	// arrange
	sut := routeGroup{
		basePath: "/tests",
		ServeMux: http.NewServeMux(),
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/tests/test1", nil)
	called := false

	// act
	sut.HandleFunc(http.MethodGet, "/test1", func(w http.ResponseWriter, r *http.Request) { called = true })
	sut.ServeHTTP(w, r)

	// assert
	if !called {
		t.Fatalf(`Handle(http.MethodGet, "/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true})) = %v want true`, called)
	}
}
