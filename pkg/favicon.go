// Package pkg provides utils to return a default favicon.
package pkg

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Favicon provides decorator to return a default favicon.
type Favicon struct {
	// raw bytes representation of the default favicon.
	data []byte
}

// NewFaviconFromFile creates a Favicon from a file.
func NewFaviconFromFile(filepath string) (*Favicon, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &Favicon{}, err
	}
	return &Favicon{data}, nil
}

// isGettingFavicon checks whether an url is requesting for a favicon.
func (f *Favicon) isGettingFavicon(url string) bool {
	if len(f.data) == 0 {
		return false
	}
	return strings.ToLower(filepath.Base(url)) == "favicon.ico"
}

// GetDefaultFavicon decorates a GetObject function to return a default favicon
// if the S3 compatible backend does not provides a favicon.
func (f *Favicon) GetDefaultFavicon(GetObject func(url string) (Resource, error)) func(url string) (Resource, error) {

	return func(url string) (Resource, error) {
		resource, error := GetObject(url)
		if error != nil && f.isGettingFavicon(url) {
			return Resource{Data: bytes.NewReader(f.data), Info: ResourceInfo{ContentType: "image/x-icon", Size: int64(len(f.data))}}, nil
		}
		return resource, error
	}
}
