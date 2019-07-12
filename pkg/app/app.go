package app

import (
	"fmt"
	"net/http"
	"time"

	core "github.com/e2fyi/minio-web/pkg/core"
	minio "github.com/e2fyi/minio-web/pkg/minio"
)

// App holds the state for the app
type App struct {
	Config Configuration
	Helper *minio.Helper
	Core
}

// NewApp creates a new App.
func NewApp() *App {
	return &App{Core: core.NewCore()}
}

// LoadConfig loads the config from both file and environment variables.
func (app *App) LoadConfig() *App {

	configuration, err := LoadConfig()
	if err != nil {
		app.Sugar.Fatal(err)
	}
	app.Config = configuration
	return app
}

// ConfigMinioHelper creates a internal helper to interact with the S3
// compatible backend.
func (app *App) ConfigMinioHelper(config MinioConfig, bucketName string, prefix string) *App {

	helper, err := minio.NewMinioHelperWithBucket(config, bucketName, prefix, 5)
	if err != nil {
		app.Sugar.Fatal(err)
	}

	// try connection for 5 times before throwing error
	app.Sugar.Infof("connecting to %s", config.Endpoint)
	for i := 1; i <= 5; i++ {
		time.Sleep(1 * time.Second)
		msg, err := helper.TestConnection()
		if err == nil {
			app.Sugar.Info(msg)
			app.Sugar.Infof("object prefix: %s", prefix)
			break
		}

		app.Sugar.Infof("attempt %d - unable to connect: %s", i, err)
		if err != nil && i >= 5 {
			app.Sugar.Fatal("unable to connect to endpoint")
		}
	}

	app.Helper = &helper
	app.ChainStatObject(helper.StatObject)
	app.ChainGetObject(helper.GetObject)

	return app
}

// StartServer creates and starts a http (or https if ssl certs are provided).
func (app *App) StartServer(config ServerConfig) {
	// flush log when app exits
	defer app.Sugar.Sync()
	// initialize core
	app.Init()
	// address to listen
	addr := fmt.Sprintf(":%d", config.Port)
	app.Sugar.Infof("Listening to port: %d", config.Port)
	// set handler
	http.HandleFunc("/", app.Handler())
	// start server
	if config.SSL.Cert == "" || config.SSL.Key == "" {
		app.Sugar.Fatal(http.ListenAndServe(addr, nil))
	} else {
		app.Sugar.Fatal(http.ListenAndServeTLS(
			addr,
			config.SSL.Cert,
			config.SSL.Key,
			nil))
	}
}
