package main

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"testing"
	"time"

	"gorest/internal/tools"
)

// wait time for server startup or shutdown
var mainWaitTime = 200

// test main() server error
// should log error
func TestMainError(t *testing.T) {
	//arrange
	defer func() {
		tools.InfoLogger.SetOutput(os.Stdin)
		tools.ErrorLogger.SetOutput(os.Stderr)
	}()
	var infBuf, errBuf bytes.Buffer
	tools.InfoLogger.SetOutput(&infBuf)
	tools.ErrorLogger.SetOutput(&errBuf)
	wantErrLog := regexp.MustCompile("http: Server closed")

	// act
	go main()
	ticker := time.NewTicker(time.Millisecond * time.Duration(mainWaitTime))
	defer ticker.Stop()
	<-ticker.C
	srv.Shutdown(context.TODO())
	<-ticker.C

	// assert
	errOutput := errBuf.String()
	if !wantErrLog.MatchString(errOutput) {
		t.Fatalf("main() = %q want %q", errOutput, wantErrLog)
	}
}

// test main() success
// should start on default port
func TestMainSuccess(t *testing.T) {
	//arrange
	defer func() {
		tools.InfoLogger.SetOutput(os.Stdin)
		tools.ErrorLogger.SetOutput(os.Stderr)
		srv.Shutdown(context.TODO())
	}()
	var infBuf, errBuf bytes.Buffer
	tools.InfoLogger.SetOutput(&infBuf)
	tools.ErrorLogger.SetOutput(&errBuf)
	wantInfLog := regexp.MustCompile("starting server on :8080")
	// act
	go main()
	ticker := time.NewTicker(time.Millisecond * time.Duration(mainWaitTime))
	defer ticker.Stop()
	<-ticker.C

	// assert
	infOutput := infBuf.String()
	if !wantInfLog.MatchString(infOutput) {
		t.Fatalf("main() = %q want %q", infOutput, wantInfLog)
	}
}

// test main() success with env
// should start on env port
func TestMainSuccessEnv(t *testing.T) {
	//arrange
	defer func() {
		tools.InfoLogger.SetOutput(os.Stdin)
		tools.InfoLogger.SetOutput(os.Stderr)
		srv.Shutdown(context.TODO())
	}()
	var infBuf, errBuf bytes.Buffer
	tools.InfoLogger.SetOutput(&infBuf)
	tools.ErrorLogger.SetOutput(&errBuf)
	os.Setenv("PORT", "1234")
	wantInfLog := regexp.MustCompile("starting server on :1234")

	// act
	go main()
	ticker := time.NewTicker(time.Millisecond * time.Duration(mainWaitTime))
	defer ticker.Stop()
	<-ticker.C

	// assert
	infOutput := infBuf.String()
	if !wantInfLog.MatchString(infOutput) {
		t.Fatalf("main() = %q want %q", infOutput, wantInfLog)
	}
}
