package ext

import (
	"fmt"
	"path"
)

// IndexHTML provides the decorator to insert a default index file to any
// requests.
type IndexHTML struct {
	filename string
}

// DefaultIndexFileExtension installs the extension where a default index file
// is queried if not provided in the url
// (e.g. http://abc instead of http://abc/efg.html).
func DefaultIndexFileExtension(filename string) Extension {
	return func(c *Core) (string, error) {
		if filename == "" {
			filename = "index.html"
		}
		decorator := GetIndexFileDecorator(filename)
		c.ApplyStatObject(decorator)
		c.ApplyGetObject(decorator)
		return fmt.Sprintf("default index file: %s", filename), nil
	}
}

// GetIndexFileDecorator installs the extension where a default index file is queried
// if not provided in the url (e.g. http://abc instead of http://abc/efg.html).
func GetIndexFileDecorator(filename string) HandlerDecorator {
	if filename == "" {
		filename = "index.html"
	}
	ext := IndexHTML{filename: filename}
	return ext.GetIndexHTML
}

// insertIfNeeded inserts the default index file into the url if needed (e.g.
// Directory request).
func (i IndexHTML) insertIfNeeded(url string) (updatedURL string, unchanged bool) {
	switch n := len(url); {
	case n == 0:
		return i.filename, false
	case n == 1 && url[0] == '/':
		return i.filename, false
	case url[n-1] == '/':
		return path.Join(url, i.filename), false
	default:
		return url, true
	}
}

// GetIndexHTML decorates a GetObject or StatObject function to insert the
// default index file if required.
func (i IndexHTML) GetIndexHTML(handler Handler) Handler {

	return func(url string) (Resource, error) {
		updatedURL, unchanged := i.insertIfNeeded(url)
		resource, error := handler(updatedURL)
		if error != nil && unchanged {
			return handler(path.Join(url, i.filename))
		}
		return resource, error
	}
}
