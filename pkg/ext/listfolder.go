package ext

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"time"

	core "github.com/e2fyi/minio-web/pkg/core"
)

const listingTemplate = `
## {{.BucketName}}{{.URL}}

- <a href="..">. .</a>
{{range .ListingItems}}
- <a href="{{.Path}}" title="{{.LastModified}}">{{.Name}}</a>
{{end}}
`

type listingItem struct {
	Name         string
	Path         string
	LastModified time.Time
}

type listing struct {
	BucketName   string
	Prefix       string
	URL          string
	ListingItems []listingItem
}

// ListFolderExt describes the extension to list objects inside a pseudo-minio folder.
type ListFolderExt struct {
	helper             *MinioHelper
	listFolder         bool
	listFolderObjects  bool
	listFolderTemplate *template.Template
}

// ListFolderExtension installs the extension to list folder objects.
func ListFolderExtension(helper *MinioHelper, listFolder bool, listFolderObjects bool) Extension {
	return func(c *Core) (string, error) {
		ext, err := NewListFolderExt(helper, listFolder, listFolderObjects)
		if err != nil {
			return "list folder: errored", err
		}
		c.ChainGetObject(ext.ListObjectsAsMarkdown)
		return "list folder: enabled", nil
	}
}

// ListFolderExtension installs the extension to list folder objects.
func (ext *ListFolderExt) decorateListFolderHandler(handler Handler) Handler {
	return core.ChainHandlers(handler, ext.ListObjectsAsMarkdown)
}

// NewListFolderExt creates a new ListFolderExt object.
func NewListFolderExt(helper *MinioHelper, listFolder bool, listFolderObjects bool) (*ListFolderExt, error) {
	if !listFolder {
		return &ListFolderExt{listFolder: listFolder}, nil
	}
	listFolderTemplate, err := template.New("listing").Parse(listingTemplate)
	if err == nil {
		return &ListFolderExt{
			helper:             helper,
			listFolder:         listFolder,
			listFolderObjects:  listFolderObjects,
			listFolderTemplate: listFolderTemplate}, nil
	}
	return &ListFolderExt{}, err
}

// ListObjectsAsMarkdown retrieves (non-recursive) objects with a specified prefix
// and rendered them as markdown Resource.
func (ext *ListFolderExt) ListObjectsAsMarkdown(url string) (Resource, error) {

	bucketName, prefix := ext.helper.GetBucketNameAndPrefix(url)
	if bucketName == "" {
		return Resource{}, errors.New("No bucket provided")
	}

	// Create a done channel to control 'ListObjectsV2' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// list folders inside folder
	isRecursive := false
	objectCh := ext.helper.Client.ListObjectsV2(bucketName, prefix, isRecursive, doneCh)
	var items []listingItem
	for info := range objectCh {
		if !ext.listFolderObjects && info.Size != 0 {
			continue
		}
		name := strings.Replace(info.Key, prefix, "", 1)
		if name[0] == '.' {
			continue
		}
		items = append(items,
			listingItem{
				Name:         name,
				Path:         name,
				LastModified: info.LastModified})
	}

	var renderedMarkdown bytes.Buffer
	err := ext.listFolderTemplate.Execute(&renderedMarkdown,
		listing{
			BucketName:   ext.helper.BucketName,
			Prefix:       prefix,
			URL:          url,
			ListingItems: items})
	if err != nil {
		return Resource{}, err
	}

	return Resource{
		Data: bytes.NewReader(renderedMarkdown.Bytes()),
		Info: ResourceInfo{
			Size:         int64(len(renderedMarkdown.Bytes())),
			ContentType:  "text/markdown",
			LastModified: time.Now()}}, nil
}
