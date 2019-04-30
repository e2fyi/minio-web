// Package pkg provides utils to start a http or https server.
package pkg

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

// SugaredLogger is alias of zap.SugaredLogger
type SugaredLogger = zap.SugaredLogger

// ServeFunction is a function that serve a resource to a http response writer.
type ServeFunction = func(http.ResponseWriter, Resource) error

// GetFunction is a function that returns a Resource based on an url.
type GetFunction = func(url string) (Resource, error)

// App holds the state for the app
type App struct {
	sugar   *SugaredLogger
	Config  Configuration
	helper  MinioHelper
	handler *Handlers
}

// NewApp creates a new App.
func NewApp() *App {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	return &App{
		sugar:   sugar,
		handler: &Handlers{Serve: DefaultServe, SetHeaders: SetDefaultHeaders}}
}

// StartServer creates and starts a http (or https if ssl certs are provided).
func (app *App) StartServer(config ServerConfig) {
	// flush log when app exits
	defer app.sugar.Sync()
	// address to listen
	addr := fmt.Sprintf(":%d", config.Port)
	app.sugar.Infof("Listening to port: %d", config.Port)
	// set handler
	http.HandleFunc("/", app.handler.Handler())
	if config.SSL.Cert == "" || config.SSL.Key == "" {
		app.sugar.Fatal(http.ListenAndServe(addr, nil))
	} else {
		app.sugar.Fatal(http.ListenAndServeTLS(
			addr,
			config.SSL.Cert,
			config.SSL.Key,
			nil))
	}
}
