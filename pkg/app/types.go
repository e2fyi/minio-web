package app

import (
	"net/http"

	"go.uber.org/zap"

	core "github.com/e2fyi/minio-web/pkg/core"
	minio "github.com/e2fyi/minio-web/pkg/minio"
)

// Core is an alias for core.Core
type Core = core.Core

// Resource is an alias for core.Resource
type Resource = core.Resource

// ResourceInfo is an alias for core.ResourceInfo
type ResourceInfo = core.ResourceInfo

// Handler is an alias for core.Handler
type Handler = core.Handler

// SugaredLogger is alias of zap.SugaredLogger
type SugaredLogger = zap.SugaredLogger

// ServeFunction is a function that serve a resource to a http response writer.
type ServeFunction = func(http.ResponseWriter, core.Resource) error

// GetFunction is a function that returns a Resource based on an url.
type GetFunction = func(url string) (core.Resource, error)

// Extension is an alias of core.Extension.
type Extension = core.Extension

// MinioConfig is an alias for core.Config
type MinioConfig = minio.Config
