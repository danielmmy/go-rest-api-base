package handlers

import (
	"bytes"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"gorest/api"
	"gorest/internal/tools"
)

// test writeResponse(w http.ResponseWriter, code int, data any, headers ...http.Header) with Write call error
// should log and return error code
func TestWriteResponseWriteError(t *testing.T) {
	// arrange
	ResetTest()
	ResponseWriterResponseError = errors.New("test error")
	w := MockResponseWriter{}
	sut := NewApp()
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrorMessage := "test error"
	wantLog := regexp.MustCompile(wantErrorMessage)
	wantCode := 1

	// act
	err := sut.writeResponse(w, wantCode, nil)

	// assert
	if err.Error() != wantErrorMessage {
		t.Fatalf("writeResponse(w, wantCode, nil) = %q want %q", err, wantErrorMessage)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("writeResponse(w, wantCode, nil) = %q want %q", errOutput, wantLog)
	}
}

// test writeResponse(w http.ResponseWriter, code int, data any, headers ...http.Header) with Data error
// should log and return error
func TestWriteResponseWithDataError(t *testing.T) {
	// arrange
	ResetTest()
	w := MockResponseWriter{}
	sut := NewApp()
	wantErrorMessage := regexp.MustCompile("json: unsupported value")
	wantErrorType := "*json.UnsupportedValueError"
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)

	// act
	err := sut.writeResponse(w, 200, math.Inf(1))

	// assert
	if reflect.TypeOf(err).String() != wantErrorType {
		t.Fatalf("writeResponse(w, 200, math.Inf(1)) = %q want %q", reflect.TypeOf(err).String(), wantErrorType)
	}

	errOutput := errBuf.String()
	if !wantErrorMessage.MatchString(errOutput) {
		t.Fatalf("writeResponse(w, 200, math.Inf(1)) = %q want %q", errOutput, wantErrorMessage)
	}

}

// test writeResponse(w http.ResponseWriter, code int, data any, headers ...http.Header) with default params
// should write code
func TestWriteResponseNullData(t *testing.T) {
	// arrange
	ResetTest()
	w := MockResponseWriter{}
	sut := NewApp()
	var wantError error = nil
	wantCode := 1

	// act
	err := sut.writeResponse(w, wantCode, nil)

	// assert
	if err != wantError {
		t.Fatalf("writeResponse(w, wantCode, nil) = %v want %v", err, wantError)
	}

	if ResponseWriterReceivedStatusCode != wantCode {
		t.Fatalf("writeResponse(w, wantCode, nil) = %d want %d", ResponseWriterReceivedStatusCode, wantCode)
	}

	if ResponseWriterReceivedBytes != nil {
		t.Fatalf("writeResponse(w, wantCode, nil) = %v want %v", ResponseWriterReceivedBytes, nil)
	}
}

// test writeResponse(w http.ResponseWriter, code int, data any, headers ...http.Header) with Custom header
// should write custom header
func TestWriteResponseCustomHeader(t *testing.T) {
	// arrange
	ResetTest()
	w := MockResponseWriter{}
	sut := NewApp()
	wantHeader := "text/plain"
	header := http.Header{
		"Content-Type": {wantHeader},
	}

	// act
	sut.writeResponse(w, 200, nil, header)

	// assert
	if ResponseWriterHeader.Values("Content-Type")[0] != wantHeader {
		t.Fatalf("writeResponse(w, 200, nil, header) = %q want %q", ResponseWriterHeader.Values("Content-Type")[0], wantHeader)
	}
}

// test writeResponse(w http.ResponseWriter, code int, data any, headers ...http.Header) with Data
// should write Data
func TestWriteResponseWithData(t *testing.T) {
	// arrange
	ResetTest()
	w := MockResponseWriter{}
	sut := NewApp()
	wantBytesMsg := `{"id":1,"owner":"owner1"}`
	fed := &api.Federation{
		Id:    1,
		Owner: "owner1",
	}

	// act
	sut.writeResponse(w, 200, fed)

	// assert
	if string(ResponseWriterReceivedBytes[:]) != wantBytesMsg {
		t.Fatalf("sut.writeResponse(w, 200, federation) = %q want %q", string(ResponseWriterReceivedBytes[:]), wantBytesMsg)
	}
}

// test writeResponse(w http.ResponseWriter, code int, data any, headers ...http.Header) with error as Data
// should write error as Data
func TestWriteResponseWithErrorAsData(t *testing.T) {
	// arrange
	ResetTest()
	w := MockResponseWriter{}
	sut := NewApp()
	wantBytesMsg := `{"msg":"test error"}`
	errData := errors.New("test error")

	// act
	sut.writeResponse(w, 400, errData)

	// assert
	if string(ResponseWriterReceivedBytes[:]) != wantBytesMsg {
		t.Fatalf("writeResponse(w, 400, errData) = %q want %q", string(ResponseWriterReceivedBytes[:]), wantBytesMsg)
	}
}

// test readJson(w http.ResponseWriter, r *http.Request, data any) error decoding data
// should log and return error
func TestReadJsonDecodeError(t *testing.T) {
	// arrange
	body := strings.NewReader(`{"id":abc,"owner":"owner1"}`)
	w := MockResponseWriter{}
	r := httptest.NewRequest("GET", "/", body)
	federation := new(api.Federation)
	sut := NewApp()
	wantLog := regexp.MustCompile(`invalid character 'a' looking for beginning of value`)
	wantErrorType := "*json.SyntaxError"
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)

	// act
	err := sut.readJson(w, r, federation)

	// assert
	if reflect.TypeOf(err).String() != wantErrorType {
		t.Fatalf("readJson(w, r, federation) = %q want %q", reflect.TypeOf(err).String(), wantErrorType)
	}

	errOutput := errBuf.String()
	if !wantLog.MatchString(errOutput) {
		t.Fatalf("readJson(w, r, federation) = %q want %q", errOutput, wantLog)
	}
}

// test readJson(w http.ResponseWriter, r *http.Request, data any) error two bodies
// should log and return error
func TestReadJsonTwoBodiesError(t *testing.T) {
	// arrange
	body := strings.NewReader(`{"id":1,"owner":"owner1"}{}`)
	w := MockResponseWriter{}
	r := httptest.NewRequest("GET", "/", body)
	federation := new(api.Federation)
	sut := NewApp()
	wantErrorMessage := "only one json body allowed"
	wantErrorType := "*errors.errorString"
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)

	// act
	err := sut.readJson(w, r, federation)

	// assert
	if reflect.TypeOf(err).String() != wantErrorType {
		t.Fatalf("readJson(w, r, federation) = %q want %q", reflect.TypeOf(err).String(), wantErrorType)
	}

	if err.Error() != wantErrorMessage {
		t.Fatalf("readJson(w, r, federation) = %q want %q", err, wantErrorMessage)
	}
}

// test readJson(w http.ResponseWriter, r *http.Request, data any) success
// should unmarshall and return nil
func TestReadJsonSuccess(t *testing.T) {
	// arrange
	body := strings.NewReader(`{"id":1,"owner":"owner1"}`)
	w := MockResponseWriter{}
	r := httptest.NewRequest("GET", "/", body)
	federation := new(api.Federation)
	sut := NewApp()
	defer func() {
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var errBuf bytes.Buffer
	tools.ErrorLogger.SetOutput(&errBuf)

	// act
	err := sut.readJson(w, r, federation)

	// assert
	if err != nil {
		t.Fatalf("readJson(w, r, federation) = %q want <nil>", reflect.TypeOf(err).String())
	}

	if federation.Id != 1 {
		t.Fatalf("readJson(w, r, federation) = %d want %d", federation.Id, 1)
	}

	if federation.Owner != "owner1" {
		t.Fatalf("readJson(w, r, federation) = %q want %q", federation.Owner, "owner1")
	}

}
