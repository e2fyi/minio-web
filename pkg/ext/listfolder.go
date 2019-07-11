package ext

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	glob "github.com/gobwas/glob"

	core "github.com/e2fyi/minio-web/pkg/core"
)

const listingTemplate = `
## {{.BucketName}}{{.URL}}

| Name | Last Modified | Size |
| --- | --- | --- |
| <a href="..">. .</a> | | |
{{range .ListingItems}} | <a href="{{.Path}}">{{.Name}}</a> | {{.LastModified}} | {{.Size}} |  
{{end}}
`

type listingItem struct {
	Name         string
	Path         string
	Size         string
	LastModified string
}

type listing struct {
	BucketName   string
	Prefix       string
	URL          string
	ListingItems []listingItem
}

// ListFolderExt describes the extension to list objects inside a pseudo-minio folder.
type ListFolderExt struct {
	pattern            glob.Glob
	helper             *MinioHelper
	listFolder         bool
	listFolderObjects  string
	listFolderTemplate *template.Template
}

// ListFolderExtension installs the extension to list folder objects.
func ListFolderExtension(helper *MinioHelper, listFolder bool, listFolderObjects string) Extension {
	return func(c *Core) (string, error) {

		ext, err := NewListFolderExt(helper, listFolder, listFolderObjects)
		if err != nil {
			return "list folder: errored", err
		}
		c.ChainGetObject(ext.ListObjectsAsMarkdown)
		return fmt.Sprintf("list folder objects: %s", listFolderObjects), nil
	}
}

// ListFolderExtension installs the extension to list folder objects.
func (ext *ListFolderExt) decorateListFolderHandler(handler Handler) Handler {
	return core.ChainHandlers(handler, ext.ListObjectsAsMarkdown)
}

// NewListFolderExt creates a new ListFolderExt object.
func NewListFolderExt(helper *MinioHelper, listFolder bool, listFolderObjects string) (*ListFolderExt, error) {
	if !listFolder {
		return &ListFolderExt{listFolder: listFolder}, nil
	}
	listFolderTemplate, err := template.New("listing").Parse(listingTemplate)

	if err == nil {
		return &ListFolderExt{
			pattern:            glob.MustCompile(listFolderObjects),
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
		// get actual filename
		name := strings.Replace(info.Key, prefix, "", 1)

		// filter hidden files
		if name[0] == '.' {
			continue
		}

		// filter objects don't match glob
		if info.Size != 0 && !ext.pattern.Match(name) {
			continue
		}

		// dun render 0 bytes
		var size string
		if info.Size == 0 {
			size = ""
		} else {
			size = humanize.Bytes(uint64(info.Size))
		}

		// dun render 0 timestamp
		var lastModified string
		if info.LastModified == (time.Time{}) {
			lastModified = ""
		} else {
			lastModified = humanize.Time(info.LastModified)
		}

		items = append(items,
			listingItem{
				Name:         name,
				Path:         name,
				Size:         size,
				LastModified: lastModified})
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
