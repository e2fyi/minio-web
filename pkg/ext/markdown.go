package ext

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
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

// RenderMarkdownExtension installs the markdown extension if a template is
// provided.
func RenderMarkdownExtension(templateFile string) Extension {
	return func(c *Core) (string, error) {
		decorator, err := getMarkdownDecorator(templateFile)
		if err != nil {
			return "markdown rendering: errored", err
		}
		c.ApplyServe(decorator)
		return "markdown rendering: enabled", nil
	}
}

// getMarkdownDecorator returns a ServeHandlerDecorator.
func getMarkdownDecorator(templateFile string) (ServeHandlerDecorator, error) {
	template, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}

	md := markdown.New(
		markdown.HTML(true),
		markdown.Tables(true),
		markdown.Linkify(true),
		markdown.Typographer(true),
		markdown.XHTMLOutput(true))

	ext := Markdown{template: template, md: md}
	return ext.RenderMarkdown, nil
}

// isMarkdown checks whether a resource content type is a markdown.
func isMarkdown(resource Resource) bool {
	return strings.Contains(strings.ToLower(resource.Info.ContentType), "markdown") || strings.HasSuffix(strings.ToLower(resource.Info.Key), ".md")
}

// RenderMarkdown decorates a Serve function to render and return a HTML
// resource from a markdown resource.
func (m Markdown) RenderMarkdown(Serve ServeHandler) ServeHandler {

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
		var buf bytes.Buffer
		err = m.template.Execute(&buf, TemplateData{Content: template.HTML(rendered)})
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(buf.Bytes())), 10))
		w.Write(buf.Bytes())
		return Serve(w, resource)
	}
}
