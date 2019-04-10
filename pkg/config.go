// Package pkg provides utils to configure the web server.
package pkg

import (
	"github.com/tkanos/gonfig"
	"log"
	"os"
)

// Configuration provides the settings to config the web server and its plugins.
type Configuration struct {
	Endpoint         string `env:"MINIO_ENDPOINT"`
	Port             int    `env:"MINIO_PORT"`
	AccessKeyID      string `env:"MINIO_ACCESS_KEY_ID"`
	SecretAccessKey  string `env:"MINIO_SECRET_ACCESS_KEY"`
	Ssl              bool   `env:"MINIO_SSL"`
	Region           string `env:"MINIO_REGION"`
	BucketName       string `env:"MINIO_BUCKET"`
	DefaultHTML      string `env:"MINIO_DEFAULT_HTML"`
	FavIcon          string `env:"MINIO_FAVICON"`
	CacheSize        int    `env:"MINIO_CACHESIZE"`
	MarkdownTemplate string `env:"MINIO_MD_TEMPLATE"`
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
func LoadConfig() Configuration {
	configuration := Configuration{}
	configFile := configFilePath()
	err := gonfig.GetConf(configFile, &configuration)
	if err != nil {
		panic(err)
	}
	log.Printf("loaded config at %s.", configFile)
	return configuration
}
