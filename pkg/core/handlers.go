package core

import (
	"io"
	"net/http"
	"strconv"
	"time"
)

// ResourceInfo describes the metadata of the resource.
type ResourceInfo struct {
	Key          string
	Size         int64
	ETag         string
	ContentType  string
	LastModified time.Time
}

// Resource represents the retrieved resource from the S3 compatible backend
// to be streamed.
type Resource struct {
	Data io.Reader
	Info ResourceInfo
	Msg  string
}

// Handlers describes how to get Resource metadata, retrieve Resource from
// S3 compatible backend, how to set the Headers, as well as how to serve
// the Resource.
type Handlers struct {
	StatObject func(url string) (Resource, error)
	GetObject  func(url string) (Resource, error)
	ListFolder func(url string) (Resource, error)
	SetHeaders func(w http.ResponseWriter, info ResourceInfo)
	Serve      func(w http.ResponseWriter, r Resource) error
	Sugared
}

// Handler is a returns a Resource with a provided url.
type Handler = func(url string) (Resource, error)

// HandlerDecorator decorates a Handler.
type HandlerDecorator = func(Handler) Handler

// ServeHandler handles the serving of a resource.
type ServeHandler = func(w http.ResponseWriter, r Resource) error

// ServeHandlerDecorator decorates a ServeHandler.
type ServeHandlerDecorator = func(ServeHandler) ServeHandler

// HeaderHandler handles the response header with the provided resource info.
type HeaderHandler = func(w http.ResponseWriter, info ResourceInfo)

// HeaderHandlerDecorator decorates a HeaderHandler.
type HeaderHandlerDecorator = func(HeaderHandler) HeaderHandler

// SetDefaultHeaders set headers for the http response.
func SetDefaultHeaders(w http.ResponseWriter, info ResourceInfo) {
	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("ETag", info.ETag)
	w.Header().Set("Last-Modified", info.LastModified.Format(time.RFC1123))
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
}

// DefaultServe serve the Resource.
func DefaultServe(w http.ResponseWriter, r Resource) error {
	if r.Data == nil {
		w.WriteHeader(404)
		return nil
	}
	_, err := io.Copy(w, r.Data)
	return err
}

// Handler sets some defaults and returns a handler function.
func (h *Handlers) Handler() func(w http.ResponseWriter, r *http.Request) {
	if h.Serve == nil {
		h.Serve = DefaultServe
	}
	if h.SetHeaders == nil {
		h.SetHeaders = SetDefaultHeaders
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "HEAD":
			h.HeadHandler(w, r)
		case "GET":
			h.GetHandler(w, r)
		default:
			w.WriteHeader(405)
		}
	}
}

// HeadHandler handles the request when method is HEAD.
func (h *Handlers) HeadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	res, err := h.StatObject(url)
	if res.Msg != "" {
		h.Sugar.Info(res.Msg)
	}
	if err != nil {
		h.Sugar.Errorf("HEAD[%s] [404]: %s", url, err.Error())
		w.WriteHeader(404)
		return
	}
	h.SetHeaders(w, res.Info)
}

// GetHandler handles the request when method is GET.
func (h *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	res, err := h.GetObject(url)
	if res.Msg != "" {
		h.Sugar.Info(res.Msg)
	}
	if err != nil {
		h.Sugar.Errorf("GET[%s] [404]: %s", url, err)
		w.WriteHeader(404)
		return
	}

	h.SetHeaders(w, res.Info)
	err = h.Serve(w, res)
	if res.Msg != "" {
		h.Sugar.Info(res.Msg)
	}
	if err != nil {
		h.Sugar.Errorf("Serve[%s] [500]: %s", url, err)
		w.Header().Set("Status-Code", "500")
		w.Write([]byte(err.Error()))
		return
	}
	h.Sugar.Infof("GET[%s] [200]: ok", res.Info.Key)
}
