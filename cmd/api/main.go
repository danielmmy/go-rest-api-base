package main

import (
	"net/http"
	"os"

	"gorest/internal/handlers"
	"gorest/internal/tools"
)

var srv = http.Server{}

func init() {

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := handlers.NewApp(handlers.WithPort(port))
	srv.Addr = app.GetAddr()
	srv.Handler = app.NewHandler()

	tools.InfoLogger.Printf("starting server on %s...\n", app.GetAddr())
	tools.ErrorLogger.Println(srv.ListenAndServe())
}
