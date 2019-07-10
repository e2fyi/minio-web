package minio

import (
	core "github.com/e2fyi/minio-web/pkg/core"
)

// Resource is an alias for core.Resource
type Resource = core.Resource

// ResourceInfo is an alias for core.ResourceInfo
type ResourceInfo = core.ResourceInfo

// Handler is an alias for core.Handler
type Handler = core.Handler

// ServeHandler is an alias for core.ServeHandler
type ServeHandler = core.ServeHandler

// minioHandler is a method that returns a Resource from a URL
type minioHandler = func(url string) (Resource, error)
