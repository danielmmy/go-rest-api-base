package handlers

import "testing"

// test NewApp() without options
// should return default app
func TestNewAppSuccessNoOption(t *testing.T) {
	// act
	app := NewApp()

	// assert
	if app.Host != "" {
		t.Fatalf(`NewApp() = %q want ""`, app.Host)
	}

	if app.Port != "" {
		t.Fatalf(`NewApp() = %q want ""`, app.Port)
	}
}

// test NewApp() with host options
// should return app with setted host
func TestNewAppSuccessHostOption(t *testing.T) {
	// act
	app := NewApp(WithHost("testHost"))

	// assert
	if app.Host != "testHost" {
		t.Fatalf(`NewApp() = %q want "testHost"`, app.Host)
	}
}

// test NewApp() with port options
// should return app with setted port
func TestNewAppSuccessPortOption(t *testing.T) {
	// act
	app := NewApp(WithPort("1234"))

	// assert
	if app.Port != "1234" {
		t.Fatalf(`NewApp() = %q want "1234"`, app.Port)
	}
}

// test GetAddr
// should return app address
func TestGetAddr(t *testing.T) {
	// arrange
	wantAddr := "testHost:1234"
	sut := NewApp(WithHost("testHost"), WithPort("1234"))

	// act
	addr := sut.GetAddr()

	// assert
	if addr != wantAddr {
		t.Fatalf("GetAddr() = %q want %q", addr, wantAddr)
	}
}
