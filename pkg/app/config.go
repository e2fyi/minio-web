package app

import (
	"os"
	"strings"

	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/env"
	"go-micro.dev/v4/config/source/file"
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

// ExtensionsConfig is used to config the extensions to install on minio-web.
type ExtensionsConfig struct {
	BucketName        string `json:"bucketname"`
	Prefix            string `json:"prefix"`
	DefaultHTML       string `json:"defaulthtml"`
	DefaultHTMLs      []string
	FavIcon           string `json:"favicon"`
	CacheSize         int    `json:"cachesize"`
	MarkdownTemplate  string `json:"markdowntemplate"`
	ListFolder        bool   `json:"listfolder"`
	ListFolderObjects string `json:"listfolderobjects"`
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
func LoadConfig() (Configuration, error) {
	// load from file
	conf, err := config.NewConfig()
	if err != nil {
		return Configuration{}, err
	}
	err = conf.Load(
		file.NewSource(
			file.WithPath(configFilePath()),
		),
		env.NewSource(),
	)
	if err != nil {
		return Configuration{}, err
	}
	var configuration Configuration
	err = conf.Scan(&configuration)
	if err != nil {
		return Configuration{}, err
	}
	configuration.Ext.DefaultHTMLs = strings.Split(
		strings.ReplaceAll(configuration.Ext.DefaultHTML, " ", ""),
		",")
	return configuration, nil
}
