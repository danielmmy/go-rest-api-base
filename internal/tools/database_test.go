package tools

import (
	"bytes"
	"math/rand"
	"os"
	"regexp"
	"testing"
)

// test NewFederationRepositoryError
// should return nil repository and error
func TestNewFederationRepositoryErro(t *testing.T) {
	// arrange
	r = rand.New(rand.NewSource(1))
	var errBuf bytes.Buffer
	ErrorLogger.SetOutput(&errBuf)
	defer func() {
		ErrorLogger.SetOutput(os.Stderr)
	}()

	wantError := regexp.MustCompile("random error")
	var wantRepo *FederationRepository = nil

	// act
	repo, err := NewFederationRepository()

	// assert
	if repo != wantRepo {
		t.Fatalf("NewFederationRepository() = %v want %v", repo, wantRepo)
	}

	if !wantError.MatchString(err.Error()) {
		t.Fatalf("NewFederationRepository() = %q want %q", err, wantError)
	}

	errOutput := errBuf.String()
	if !wantError.MatchString(errOutput) {
		t.Fatalf("NewFederationRepository() = %q want %q", errOutput, wantError)
	}
}

// test NewFederationRepository success
// should return *FederationRepository
func TestNewFederationRepositorySucess(t *testing.T) {
	// arrange
	var wantError error = nil
	var errBuf bytes.Buffer
	ErrorLogger.SetOutput(&errBuf)
	defer func() {
		ErrorLogger.SetOutput(os.Stderr)
	}()

	// act
	repo, err := NewFederationRepository()

	// assert
	if repo == nil {
		t.Fatal(`NewFederationRepository() = <nil> want "*FederationRepository"`)
	}

	if err != wantError {
		t.Fatalf("NewFederationRepository() = %v want %v", err, wantError)
	}

	errOutput := errBuf.String()
	if errOutput != "" {
		t.Fatalf(`NewFederationRepository() = %q want ""`, errOutput)
	}
}
