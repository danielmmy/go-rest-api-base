package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"gorest/internal/tools"
)

// test Authorize call without token
// should log error
func TestAuthorizeNilToken(t *testing.T) {
	// arrange
	want := regexp.MustCompile("unauthorized access attempt from")
	var errBuf bytes.Buffer
	var infBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	tools.InfoLogger.SetOutput(&infBuf)
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
		tools.InfoLogger.SetOutput(os.Stdout)
	}()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.InfoLogger.Println("called")
	})

	sut := Authorize(nextHandler)
	req := httptest.NewRequest("GET", "http://test", nil)

	// act
	sut.ServeHTTP(httptest.NewRecorder(), req)
	errOutput := errBuf.String()
	infOutput := infBuf.String()

	// assert
	if !want.MatchString(errOutput) {
		t.Fatalf("Authorize(nextHandler) = %q want %q", errOutput, want)
	}

	if infOutput != "" {
		t.Fatalf(`Authorize(nextHandler) = %q want ""`, infOutput)
	}
}

// test Authorize call with bad token
// should log error
func TestAuthorizeBadToken(t *testing.T) {
	// arrange
	want := regexp.MustCompile("unauthorized access attempt from")
	var errBuf bytes.Buffer
	var infBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	tools.InfoLogger.SetOutput(&infBuf)
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
		tools.InfoLogger.SetOutput(os.Stdout)
	}()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.InfoLogger.Println("called")
	})

	sut := Authorize(nextHandler)
	req := httptest.NewRequest("GET", "http://test", nil)
	req.Header.Set("Authorization", "badToken")

	// act
	sut.ServeHTTP(httptest.NewRecorder(), req)
	errOutput := errBuf.String()
	infOutput := infBuf.String()

	// assert
	if !want.MatchString(errOutput) {
		t.Fatalf("Authorize(nextHandler) = %q want %q", errOutput, want)
	}

	if infOutput != "" {
		t.Fatalf(`Authorize(nextHandler) = %q want ""`, infOutput)
	}
}

// test Authorize call with good token
// should call next
func TestAuthorizeSuccess(t *testing.T) {
	// arrange
	want := regexp.MustCompile("called")
	var errBuf bytes.Buffer
	var infBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	tools.InfoLogger.SetOutput(&infBuf)
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
		tools.InfoLogger.SetOutput(os.Stdout)
	}()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.InfoLogger.Println("called")
	})

	sut := Authorize(nextHandler)
	req := httptest.NewRequest("GET", "http://test", nil)
	req.Header.Set("Authorization", "123456")

	// act
	sut.ServeHTTP(httptest.NewRecorder(), req)
	errOutput := errBuf.String()
	infOutput := infBuf.String()

	// assert
	if errOutput != "" {
		t.Fatalf(`Authorize(nextHandler) = %q want ""`, errOutput)
	}

	if !want.MatchString(infOutput) {
		t.Fatalf("Authorize(nextHandler) = %q want %q", infOutput, want)
	}
}
