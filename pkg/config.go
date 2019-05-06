// Package pkg provides utils to configure the web server.
package pkg

import (
	"os"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
)

// Configuration is global configuration object.
type Configuration struct {
	Server ServerConfig     `json:"server"`
	Minio  MinioConfig      `json:"minio"`
	Ext    ExtensionsConfig `json:"ext"`
}

// ServerConfig is used to initialize the http server.
type ServerConfig struct {
	Port int       `json:"port"`
	SSL  SSLConfig `json:"ssl"`
}

// SSLConfig is used to config a https/http2 server.
type SSLConfig struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

// MinioConfig is used to create a minio client.
type MinioConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
	Secure    bool   `json:"secure"`
	Region    string `json:"region"`
}

// ExtensionsConfig is used to config the extensions to install on minio-web.
type ExtensionsConfig struct {
	BucketName       string `json:"bucketname"`
	DefaultHTML      string `json:"defaulthtml"`
	FavIcon          string `json:"favicon"`
	CacheSize        int    `json:"cachesize"`
	MarkdownTemplate string `json:"markdowntemplate"`
}

// configFilePath returns the location of the config file.
func configFilePath() string {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		return "configs/config.json"
	}
	return configFilePath
}

// LoadConfig loads the config from both file and environment variables.
func (app *App) LoadConfig() *App {
	// load from file
	conf := config.NewConfig()
	err := conf.Load(
		file.NewSource(
			file.WithPath(configFilePath()),
		),
		env.NewSource(),
	)
	var configuration Configuration
	err = conf.Scan(&configuration)
	if err != nil {
		app.sugar.Fatal(err)
	}
	app.Config = configuration
	return app
}
