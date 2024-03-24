package handlers

import "net/http"

type MockResponseWriter struct {
}

var ResponseWriterHeader = http.Header{}
var ResponseWriterResponseCode int
var ResponseWriterResponseError error
var ResponseWriterReceivedBytes []byte
var ResponseWriterReceivedStatusCode = 0

func ResetTest() {
	ResponseWriterHeader = http.Header{}
	ResponseWriterResponseCode = 0
	ResponseWriterResponseError = nil
	ResponseWriterReceivedBytes = nil
	ResponseWriterReceivedStatusCode = 0
}

func (m MockResponseWriter) Header() http.Header {
	return ResponseWriterHeader
}
func (m MockResponseWriter) WriteBytes(receivedBytes []byte) (int, error) {
	ResponseWriterReceivedBytes = receivedBytes
	return ResponseWriterResponseCode, ResponseWriterResponseError
}

func (m MockResponseWriter) Write(receivedBytes []byte) (int, error) {
	ResponseWriterReceivedBytes = receivedBytes
	return ResponseWriterResponseCode, ResponseWriterResponseError
}

func (m MockResponseWriter) WriteHeader(statusCode int) {
	ResponseWriterReceivedStatusCode = statusCode
}
