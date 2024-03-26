package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"gorest/api"
	"gorest/internal/tools"
)

// test addFederation(w http.ResponseWriter, r *http.Request) read error
// should call readJson, log and respond error
func TestAddFederationReadError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		called = true
		return errors.New("test error")
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	r := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "test error"

	// act
	sut.addFederation(w, r)

	// assert
	if !called {
		t.Fatalf("addFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("addFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("addFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test addFederation(w http.ResponseWriter, r *http.Request) Repository connection error
// should call repostory, log and respond error
func TestAddFederationRepoError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		called = true
		return nil, errors.New("test error")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "internal server error"

	// act
	sut.addFederation(w, r)

	// assert
	if !called {
		t.Fatalf("addFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("addFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("addFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test addFederation(w http.ResponseWriter, r *http.Request) respository insert error
// should respond error
func TestAddFederationInsertError(t *testing.T) {
	// arange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		called = true
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		FederationRepositoryMockReturnCode = 400
		FederationRepositoryMockReturnError = errors.New("test error")
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := "test error"

	// act
	sut.addFederation(w, r)

	// assert
	if !called {
		t.Fatalf("addFederation(w, r) = %v want %v", called, true)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("addFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}

	errOutput := errBuf.String()
	if errOutput != "" {
		t.Fatalf(`addFederation(w, r) = %q want ""`, errOutput)
	}
}

// test addFederation(w http.ResponseWriter, r *http.Request) respository insert error and respond error
// should call respond and log error
func TestAddFederationInsertErrorRespondError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		called = true
		receivedError = data.(error)
		return errors.New("test error 2")
	}
	repository = func() (*tools.FederationRepository, error) {
		FederationRepositoryMockReturnCode = 400
		FederationRepositoryMockReturnError = errors.New("test error")
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := "test error"
	wantLog := regexp.MustCompile("test error 2")

	// act
	sut.addFederation(w, r)

	// assert
	if !called {
		t.Fatalf("addFederation(w, r) = %v want %v", called, true)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("addFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("addFederation(w, r) = %q want %q", errOutput, wantLog)
	}
}

// test addFederation(w http.ResponseWriter, r *http.Request) insert success
// should insert federation and respond created
func TestAddFederationSuccess(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedData any
	receivedCode := 0
	federation := api.Federation{
		Id:    1,
		Owner: "Test",
	}
	readJsonAlias = func(_ *App, _ http.ResponseWriter, _ *http.Request, data any) error {
		u := data.(*api.Federation)
		u.Id = 1
		u.Owner = federation.Owner
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, code int, data any, _ ...http.Header) error {
		called = true
		receivedCode = code
		receivedData = data
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnCode = 201
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	wantCode := 201

	// act
	sut.addFederation(w, r)

	// assert
	if !called {
		t.Fatalf("addFederation(w, r) = %v want %v", called, true)
	}

	if *FederationRepositoryMockReturnReceivedFed != federation {
		t.Fatalf("addFederation(w, r) = %v want %v", *FederationRepositoryMockReturnReceivedFed, federation)
	}

	if receivedCode != wantCode {
		t.Fatalf("addFederation(w, r) = %d want %d", receivedCode, wantCode)
	}

	if receivedData != nil {
		t.Fatalf("addFederation(w, r) = %v want <nil>", receivedData)
	}
}

// test getFederation(w http.ResponseWriter, r *http.Request) Repository connection error
// should call repostory, log and respond error
func TestGetFederationBadUrlParam(t *testing.T) {
	// arrange
	sut := NewApp()
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/federations/abc", nil)
	r.SetPathValue("id", "abc")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := `strconv.Atoi: parsing "abc": invalid syntax`
	wantLog := regexp.MustCompile(wantErrorMessage)

	// act
	sut.GetFederation(w, r)

	// assert
	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("getFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("getFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test getFederation(w http.ResponseWriter, r *http.Request) Repository connection error
// should call repostory, log and respond error
func TestGetFederationRepoError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		called = true
		return nil, errors.New("test error")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/federations/1", nil)
	r.SetPathValue("id", "1")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "internal server error"

	// act
	sut.GetFederation(w, r)

	// assert
	if !called {
		t.Fatalf("getFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("getFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("getFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test getFederation(w http.ResponseWriter, r *http.Request) Repository nil return response error
// should call repostory, and log error
func TestGetFederationRepoNilRespondError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return errors.New("test error")
	}
	repository = func() (*tools.FederationRepository, error) {
		called = true
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/federations/3", nil)
	r.SetPathValue("id", "3")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "federation 3 not found"

	// act
	sut.GetFederation(w, r)

	// assert
	if !called {
		t.Fatalf("getFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("getFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("getFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test getFederation(w http.ResponseWriter, r *http.Request) Repository nil return
// should call repostory, and return not found response
func TestGetFederationRepoNil(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		called = true
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/federations/3", nil)
	r.SetPathValue("id", "3")
	wantErrorMessage := "federation 3 not found"

	// act
	sut.GetFederation(w, r)

	// assert
	if !called {
		t.Fatalf("getFederation(w, r) = %v want %v", called, true)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("getFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test getFederation(w http.ResponseWriter, r *http.Request) Error reponding
// should call repostory, and log error
func TestGetFederationErrorResponding(t *testing.T) {
	// arrange
	sut := NewApp()
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, _ any, _ ...http.Header) error {
		return errors.New("test error")
	}
	repository = func() (*tools.FederationRepository, error) {
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/federations/1", nil)
	r.SetPathValue("id", "1")
	wantLog := regexp.MustCompile("test error")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)

	// act
	sut.GetFederation(w, r)

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("getFederation(w, r) = %q want %q", errOutput, wantLog)
	}
}

// test getFederation(w http.ResponseWriter, r *http.Request) success
// should return federation
func TestGetFederationSuccess(t *testing.T) {
	// arrange
	sut := NewApp()
	receivedFederation := new(api.Federation)
	wantFederation := api.Federation{
		Id:    1,
		Owner: "Owner 1",
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedFederation = data.(*api.Federation)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/federations/1", nil)
	r.SetPathValue("id", "1")

	// act
	sut.GetFederation(w, r)

	if *receivedFederation != wantFederation {
		t.Fatalf("getFederation(w, r) = %v want %v", *receivedFederation, wantFederation)
	}
}

// test getFederations(w http.ResponseWriter, r *http.Request) Repository connection error
// should call repostory, log and respond error
func TestGetFederationsRepoError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		called = true
		return nil, errors.New("test error")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "internal server error"

	// act
	sut.getFederations(w, r)

	// assert
	if !called {
		t.Fatalf("getFederations(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("getFederations(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("getFederations(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test getFederations(w http.ResponseWriter, r *http.Request) Error reponding
// should call repostory, and log error
func TestGetFederationsErrorResponding(t *testing.T) {
	// arrange
	sut := NewApp()
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, _ any, _ ...http.Header) error {
		return errors.New("test error")
	}
	repository = func() (*tools.FederationRepository, error) {
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	body := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", body)
	wantLog := regexp.MustCompile("test error")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)

	// act
	sut.getFederations(w, r)

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("getFederations(w, r) = %q want %q", errOutput, wantLog)
	}
}

// test getFederations(w http.ResponseWriter, r *http.Request) success
// should return federations
func TestGetFederationsSuccess(t *testing.T) {
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	// arrange
	sut := NewApp()
	var receivedFederations []*api.Federation
	wantFederations := []api.Federation{
		{Id: 1, Owner: "Owner 1"},
		{Id: 2, Owner: "Owner 2"},
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedFederations = data.([]*api.Federation)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	// act
	sut.getFederations(w, r)

	if len(receivedFederations) != len(wantFederations) {
		t.Fatalf("getFederations(w, r) = %d want %d", len(receivedFederations), len(wantFederations))
	}

	for i, v := range wantFederations {
		if *receivedFederations[i] != v {
			t.Fatalf("getFederations(w, r) = %v want %v", *receivedFederations[i], v)
		}
	}
}

// test updateFederation(w http.ResponseWriter, r *http.Request) read error
// should call readJson, log and respond error
func TestUpdateFederationReadError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		called = true
		return errors.New("test error")
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	r := httptest.NewRequest("PUT", "/federations/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "test error"

	// act
	sut.updateFederation(w, r)

	// assert
	if !called {
		t.Fatalf("updateFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("updateFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("updateFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test updateFederation(w http.ResponseWriter, r *http.Request) url param parse error
// should call readJson, log and respond url param error
func TestUpdateFederationUrlParamError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		called = true
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	r := httptest.NewRequest("PUT", "/federations/abc", nil)
	r.SetPathValue("id", "abc")
	w := httptest.NewRecorder()
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := `strconv.Atoi: parsing "abc": invalid syntax`
	wantLog := regexp.MustCompile(wantErrorMessage)

	// act
	sut.updateFederation(w, r)

	// assert
	if !called {
		t.Fatalf("updateFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("updateFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("updateFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test updateFederation(w http.ResponseWriter, r *http.Request) Repository connection error
// should call repostory, log and respond error
func TestUpdateFederationRepoError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		called = true
		return nil, errors.New("test error")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/federations/1", nil)
	r.SetPathValue("id", "1")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "internal server error"

	// act
	sut.updateFederation(w, r)

	// assert
	if !called {
		t.Fatalf("updateFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("updateFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("updateFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test updateFederation(w http.ResponseWriter, r *http.Request) respository update error
// should respond error
func TestUpdateFederationUpdateError(t *testing.T) {
	// arange
	sut := NewApp()
	called := false
	var receivedError error
	readJsonAlias = func(*App, http.ResponseWriter, *http.Request, any) error {
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		called = true
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		FederationRepositoryMockReturnCode = 404
		FederationRepositoryMockReturnError = errors.New("test error")
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/federations/1", nil)
	r.SetPathValue("id", "1")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := "test error"

	// act
	sut.updateFederation(w, r)

	// assert
	if !called {
		t.Fatalf("updateFederation(w, r) = %v want %v", called, true)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("updateFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}

	errOutput := errBuf.String()
	if errOutput != "" {
		t.Fatalf(`updateFederation(w, r) = %q want ""`, errOutput)
	}
}

// test updateFederation(w http.ResponseWriter, r *http.Request) respository update error and respond error
// should call respond and log error
func TestUpdateFederationUpdateErrorRespondError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	federation := api.Federation{
		Id:    1,
		Owner: "Test",
	}
	readJsonAlias = func(_ *App, _ http.ResponseWriter, _ *http.Request, data any) error {
		u := data.(*api.Federation)
		u.Id = 1
		u.Owner = federation.Owner
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		called = true
		receivedError = data.(error)
		return errors.New("test error 2")
	}
	repository = func() (*tools.FederationRepository, error) {
		FederationRepositoryMockReturnCode = 404
		FederationRepositoryMockReturnError = errors.New("test error")
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/federations/1", nil)
	r.SetPathValue("id", "1")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := "test error"
	wantLog := regexp.MustCompile("test error 2")

	// act
	sut.updateFederation(w, r)

	// assert
	if !called {
		t.Fatalf("updateFederation(w, r) = %v want %v", called, true)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("updateFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("updateFederation(w, r) = %q want %q", errOutput, wantLog)
	}
}

// test updateFederation(w http.ResponseWriter, r *http.Request) update success
// should call update federation and respond success
func TestUpdateFederationSuccess(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedData any
	receivedCode := 0
	federation := api.Federation{
		Id:    3,
		Owner: "Test",
	}
	readJsonAlias = func(_ *App, _ http.ResponseWriter, _ *http.Request, data any) error {
		u := data.(*api.Federation)
		u.Id = federation.Id
		u.Owner = federation.Owner
		return nil
	}
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, code int, data any, _ ...http.Header) error {
		called = true
		receivedCode = code
		receivedData = data
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnCode = 200
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/federations/1", nil)
	r.SetPathValue("id", "1")
	wantCode := 200

	// act
	sut.updateFederation(w, r)

	// assert
	if !called {
		t.Fatalf("updateFederation(w, r) = %v want %v", called, true)
	}

	if receivedCode != wantCode {
		t.Fatalf("updateFederation(w, r) = %d want %d", receivedCode, wantCode)
	}

	if receivedData != nil {
		t.Fatalf("updateFederation(w, r) = %v want <nil>", receivedData)
	}
}

// test deleteFederation(w http.ResponseWriter, r *http.Request) Url param parse error
// should call log and respond error
func TestDeleteFederationUrlParamError(t *testing.T) {
	// arrange
	sut := NewApp()
	var receivedError error
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		return nil, errors.New("test error")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/federations/abc", nil)
	r.SetPathValue("id", "abc")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := `strconv.Atoi: parsing "abc": invalid syntax`
	wantLog := regexp.MustCompile(wantErrorMessage)

	// act
	sut.deleteFederation(w, r)

	// assert
	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("deleteFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("deleteFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test deleteFederation(w http.ResponseWriter, r *http.Request) Repository connection error
// should call repostory, log and respond error
func TestDeleteFederationRepoError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		called = true
		return nil, errors.New("test error")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/federations/3", nil)
	r.SetPathValue("id", "3")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantLog := regexp.MustCompile("test error")
	wantErrorMessage := "internal server error"

	// act
	sut.deleteFederation(w, r)

	// assert
	if !called {
		t.Fatalf("deleteFederation(w, r) = %v want %v", called, true)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("deleteFederation(w, r) = %q want %q", errOutput, wantLog)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("deleteFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}
}

// test deleteFederation(w http.ResponseWriter, r *http.Request) respository delete error
// should respond error
func TestDeleteFederationDeleteError(t *testing.T) {
	// arange
	sut := NewApp()
	called := false
	var receivedError error
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		called = true
		receivedError = data.(error)
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		FederationRepositoryMockReturnCode = 404
		FederationRepositoryMockReturnError = errors.New("test error")
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/federations/1", nil)
	r.SetPathValue("id", "1")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := "test error"

	// act
	sut.deleteFederation(w, r)

	// assert
	if !called {
		t.Fatalf("deleteFederation(w, r) = %v want %v", called, true)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("deleteFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}

	errOutput := errBuf.String()
	if errOutput != "" {
		t.Fatalf(`deleteFederation(w, r) = %q want ""`, errOutput)
	}
}

// test deleteFederation(w http.ResponseWriter, r *http.Request) respository delete error and respond error
// should call respond and log error
func TestDeleteFederationDeleteErrorRespondError(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedError error
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, _ int, data any, _ ...http.Header) error {
		called = true
		receivedError = data.(error)
		return errors.New("test error 2")
	}
	repository = func() (*tools.FederationRepository, error) {
		FederationRepositoryMockReturnCode = 404
		FederationRepositoryMockReturnError = errors.New("test error")
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/federations/1", nil)
	r.SetPathValue("id", "1")
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := "test error"
	wantLog := regexp.MustCompile("test error 2")

	// act
	sut.deleteFederation(w, r)

	// assert
	if !called {
		t.Fatalf("deleteFederation(w, r) = %v want %v", called, true)
	}

	if receivedError.Error() != wantErrorMessage {
		t.Fatalf("deleteFederation(w, r) = %q want %q", receivedError, wantErrorMessage)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("deleteFederation(w, r) = %q want %q", errOutput, wantLog)
	}
}

// test delete Federation(w http.ResponseWriter, r *http.Request) delete success
// should call delete federation and respond success
func TestDeleteFederationSuccess(t *testing.T) {
	// arrange
	sut := NewApp()
	called := false
	var receivedData any
	receivedCode := 0
	writeResponseAlias = func(_ *App, _ http.ResponseWriter, code int, data any, _ ...http.Header) error {
		called = true
		receivedCode = code
		receivedData = data
		return nil
	}
	repository = func() (*tools.FederationRepository, error) {
		ResetFederationRepositoryMock()
		FederationRepositoryMockReturnCode = 200
		FederationRepositoryMockReturnError = nil
		return NewFederationRepositoryMock(), nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/federations/1", nil)
	r.SetPathValue("id", "1")
	wantCode := 200

	// act
	sut.deleteFederation(w, r)

	// assert
	if !called {
		t.Fatalf("deleteFederation(w, r) = %v want %v", called, true)
	}

	if receivedCode != wantCode {
		t.Fatalf("deleteFederation(w, r) = %d want %d", receivedCode, wantCode)
	}

	if receivedData != nil {
		t.Fatalf("deleteFederation(w, r) = %v want <nil>", receivedData)
	}
}
