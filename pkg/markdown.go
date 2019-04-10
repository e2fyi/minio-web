// Package pkg provides utils to render a Markdown resource as HTML.
package pkg

import (
	"bytes"
	"gitlab.com/golang-commonmark/markdown"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

// Markdown provides the decorator to serve markdowns as HMTL.
type Markdown struct {
	template *template.Template
	md       *markdown.Markdown
}

// TemplateData provides the view to the HTML template.
type TemplateData struct {
	Content template.HTML
}

// NewMarkdown creates a new Markdown from a HTML template file.
func NewMarkdown(templateFile string) (Markdown, error) {
	template, err := template.ParseFiles(templateFile)
	if err != nil {
		return Markdown{}, err
	}

	md := markdown.New(
		markdown.HTML(true),
		markdown.Tables(true),
		markdown.Linkify(true),
		markdown.Typographer(true),
		markdown.XHTMLOutput(true))

	return Markdown{template: template, md: md}, nil
}

// isMarkdown checks whether a resource content type is a markdown.
func isMarkdown(resource Resource) bool {
	return strings.Contains(strings.ToLower(resource.Info.ContentType), "markdown")
}

// RenderMarkdown decorates a Serve function to render and return a HTML
// resource from a markdown resource.
func (m Markdown) RenderMarkdown(Serve func(http.ResponseWriter, Resource) error) func(http.ResponseWriter, Resource) error {

	return func(w http.ResponseWriter, resource Resource) error {
		if !isMarkdown(resource) {
			return Serve(w, resource)
		}

		content, err := ioutil.ReadAll(resource.Data)
		resource.Data = bytes.NewReader(content)
		if err != nil {
			return Serve(w, resource)
		}

		rendered := m.md.RenderToString(content)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = m.template.Execute(w, TemplateData{Content: template.HTML(rendered)})
		if err != nil {
			return Serve(w, resource)
		}
		return err
	}
}
