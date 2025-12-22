package front

import (
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

type Renderable interface {
	GetMarkdown() string
	GetHTML() *string
	WriteHTML(string) error
}

func RenderHTML(r Renderable) (string, error) {
	var html string

	// Check if we've already rendered it
	rendered := r.GetHTML()
	if rendered != nil {
		html = *rendered
		return html, nil
	}

	// Render Markdown to HTML
	html = string(markdown.ToHTML([]byte(r.GetMarkdown()), nil, nil))
	html = bluemonday.UGCPolicy().Sanitize(html)

	// Cache rendered HTML for future reads
	err := r.WriteHTML(html)
	return html, err
}
