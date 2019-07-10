package ext

import (
	core "github.com/e2fyi/minio-web/pkg/core"
	minio "github.com/e2fyi/minio-web/pkg/minio"
)

// Core is an alias for core.Core
type Core = core.Core

// Resource is an alias for core.Resource
type Resource = core.Resource

// ResourceInfo is an alias for core.ResourceInfo
type ResourceInfo = core.ResourceInfo

// Config is an alias for core.MinioConfig
type Config = minio.Config

// Handler is an alias for core.Handler
type Handler = core.Handler

// HandlerDecorator is an alias for core.HandlerDecorator
type HandlerDecorator = core.HandlerDecorator

// ServeHandler is an alias for core.ServeHandler
type ServeHandler = core.ServeHandler

// ServeHandlerDecorator is an alias for core.ServeHandlerDecorator
type ServeHandlerDecorator = core.ServeHandlerDecorator

// Extension is an alias for core.Extension
type Extension = core.Extension

// MinioHelper is an alias for minio.MinioHelper
type MinioHelper = minio.Helper
