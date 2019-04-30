// Package pkg provides utils to return a default index file if not provided
// in the url request.
package pkg

import (
	"path"
)

// IndexHTML provides the decorator to insert a default index file to any
// requests.
type IndexHTML struct {
	filename string
}

// InsertIndexFile installs the extension where a default index file is queried
// if not provided in the url (e.g. http://abc instead of http://abc/efg.html).
func (app *App) InsertIndexFile(filename string) *App {
	if filename == "" {
		filename = "index.html"
	}
	ext := IndexHTML{filename: filename}
	app.handler.StatObject = ext.GetIndexHTML(app.handler.StatObject)
	app.handler.GetObject = ext.GetIndexHTML(app.handler.GetObject)
	app.sugar.Infof("default index file: %s", filename)
	return app
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
func (i IndexHTML) GetIndexHTML(GetObjectOrStat func(url string) (Resource, error)) func(url string) (Resource, error) {

	return func(url string) (Resource, error) {
		updatedURL, unchanged := i.insertIfNeeded(url)
		resource, error := GetObjectOrStat(updatedURL)
		if error != nil && unchanged {
			return GetObjectOrStat(path.Join(url, i.filename))
		}
		return resource, error
	}
}
