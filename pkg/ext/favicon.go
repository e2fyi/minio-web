package ext

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Favicon provides decorator to return a default favicon.
type Favicon struct {
	// raw bytes representation of the default favicon.
	data []byte
}

// DefaultFaviconExtension returns a default favicon if backend does not have one.
func DefaultFaviconExtension(filepath string) Extension {
	return func(c *Core) (string, error) {
		handler, err := newFaviconHandler(filepath)
		if err != nil {
			return "default favicon: errored", err
		}
		c.ChainGetObject(handler)
		return "default favicon: enabled", nil
	}
}

// newFaviconHandler returns favicon handler.
func newFaviconHandler(filepath string) (Handler, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	favicon := &Favicon{data}
	return favicon.Handler, nil
}

// isGettingFavicon checks whether an url is requesting for a favicon.
func (f *Favicon) isGettingFavicon(url string) bool {
	if len(f.data) == 0 {
		return false
	}
	return strings.ToLower(filepath.Base(url)) == "favicon.ico"
}

// Handler returns a default favicon for a favicon request.
func (f *Favicon) Handler(url string) (Resource, error) {
	if f.isGettingFavicon(url) {
		return Resource{
			Data: bytes.NewReader(f.data),
			Info: ResourceInfo{
				ContentType: "image/x-icon",
				Size:        int64(len(f.data))}}, nil
	}
	return Resource{}, errors.New("Not a favicon")
}
