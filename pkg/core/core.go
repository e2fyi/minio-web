package core

import (
	"go.uber.org/zap"
)

// Core describes the state of application
type Core struct {
	Sugar    *zap.SugaredLogger
	handlers *HandlersToUse
	Handlers
}

// HandlersToUse describes the types of handlers.
type HandlersToUse struct {
	StatObjects []Handler
	GetObjects  []Handler
}

// Extension is a function that decorates a Core struct.
type Extension = func(*Core) (string, error)

// NewCore creates a new Core.
func NewCore() Core {
	sugar := zap.NewExample().Sugar()
	return Core{
		Sugar: sugar,
		Handlers: Handlers{
			Serve:      DefaultServe,
			SetHeaders: SetDefaultHeaders},
		handlers: &HandlersToUse{
			StatObjects: []Handler{},
			GetObjects:  []Handler{}}}
}

// ChainStatObject applies a StatObject handler.
func (c *Core) ChainStatObject(handler Handler) *Core {
	c.handlers.StatObjects = append(c.handlers.StatObjects, handler)
	return c
}

// ChainGetObject applies a GetObject handler.
func (c *Core) ChainGetObject(handler Handler) *Core {
	c.handlers.GetObjects = append(c.handlers.GetObjects, handler)
	return c
}

// ApplyGetObject decorate the current GetObject handler.
func (c *Core) ApplyGetObject(decorator HandlerDecorator) *Core {
	// collapse all chains if a decorator is required
	getObject := ChainHandlers(c.handlers.GetObjects...)
	getObject = decorator(getObject)
	c.handlers.GetObjects = []Handler{getObject}
	return c
}

// ApplyStatObject decorate the current StatObject handler.
func (c *Core) ApplyStatObject(decorator HandlerDecorator) *Core {
	// collapse all chains if a decorator is required
	statObject := ChainHandlers(c.handlers.StatObjects...)
	statObject = decorator(statObject)
	c.handlers.StatObjects = []Handler{statObject}
	return c
}

// ApplyServe decorate the current Serve handler.
func (c *Core) ApplyServe(decorator ServeHandlerDecorator) *Core {
	c.Serve = decorator(c.Serve)
	return c
}

// ApplyHeader decorate the current Header handler.
func (c *Core) ApplyHeader(decorator HeaderHandlerDecorator) *Core {
	c.SetHeaders = decorator(c.SetHeaders)
	return c
}

// ApplyExtension applies an extension on core state.
func (c *Core) ApplyExtension(ext Extension) *Core {
	msg, err := ext(c)
	c.Sugar.Info(msg)
	if err != nil {
		c.Sugar.Fatal(err)
	}
	return c
}

// Init initialize the core state.
func (c *Core) Init() *Core {
	c.GetObject = ChainHandlers(c.handlers.GetObjects...)
	c.StatObject = ChainHandlers(c.handlers.StatObjects...)
	return c
}

// ChainHandlers returns a new handler that chains handlers sequetially
func ChainHandlers(handlers ...Handler) Handler {
	if len(handlers) == 0 {
		return nil
	} else if len(handlers) == 1 {
		return handlers[0]
	}
	return func(url string) (Resource, error) {
		var err error
		res := Resource{}

		for _, handler := range handlers {
			res, err = handler(url)
			if err == nil {
				return res, nil
			}
		}

		return res, err
	}
}
