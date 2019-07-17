package ext

import (
	"fmt"
	"path"
)

// IndexHTML provides the decorator to insert a default index file to any
// requests.
type IndexHTML struct {
	filenames []string
}

// DefaultIndexFileExtension installs the extension where a default index file
// is queried if not provided in the url
// (e.g. http://abc instead of http://abc/efg.html).
func DefaultIndexFileExtension(filenames ...string) Extension {
	return func(c *Core) (string, error) {
		if len(filenames) == 0 {
			filenames = []string{"index.html", "README.md", "index.htm"}
		}
		decorator := GetIndexFileDecorator(filenames)
		c.ApplyStatObject(decorator)
		c.ApplyGetObject(decorator)
		return fmt.Sprintf("default index files: %s", filenames), nil
	}
}

// GetIndexFileDecorator installs the extension where a default index file is queried
// if not provided in the url (e.g. http://abc instead of http://abc/efg.html).
func GetIndexFileDecorator(filenames []string) HandlerDecorator {
	ext := IndexHTML{filenames: filenames}
	return ext.GetIndexHTML
}

// insertIfNeeded inserts the default index file into the url if needed (e.g.
// Directory request).
func (i IndexHTML) getPotentialUrls(url string) []string {
	var urlCandidates = []string{}
	switch n := len(url); {
	case n == 0:
	case n == 1 && url[0] == '/':
		for _, filename := range i.filenames {
			urlCandidates = append(urlCandidates, filename)
		}
	case url[n-1] == '/':
		for _, filename := range i.filenames {
			urlCandidates = append(urlCandidates, path.Join(url, filename))
		}
	default:
		urlCandidates = append(urlCandidates, url)
		for _, filename := range i.filenames {
			urlCandidates = append(urlCandidates, path.Join(url, filename))
		}
	}
	return urlCandidates
}

// GetIndexHTML decorates a GetObject or StatObject function to insert the
// default index file if required.
func (i IndexHTML) GetIndexHTML(handler Handler) Handler {

	return func(url string) (Resource, error) {
		var resource Resource
		var error error

		urls := i.getPotentialUrls(url)

		for _, urlCandidate := range urls {
			resource, error = handler(urlCandidate)
			if error == nil {
				return resource, nil
			}
		}
		return Resource{}, error
	}
}
