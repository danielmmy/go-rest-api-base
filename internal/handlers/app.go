package handlers

import (
	"errors"
	"fmt"
)

type appOpts struct {
	Host string
	Port string
}

type App struct {
	*appOpts
}

type appConfigFunc func(*appOpts)

var errInternalServerError = errors.New("internal server error")

func NewApp(configs ...appConfigFunc) *App {
	o := appOpts{}
	for _, fn := range configs {
		fn(&o)
	}

	return &App{&o}
}

func WithHost(host string) appConfigFunc {
	return func(o *appOpts) {
		o.Host = host
	}
}

func WithPort(port string) appConfigFunc {
	return func(o *appOpts) {
		o.Port = port
	}
}

func (app *App) GetAddr() string {
	return fmt.Sprintf("%s:%s", app.Host, app.Port)
}
