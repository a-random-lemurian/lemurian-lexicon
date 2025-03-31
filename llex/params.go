package llex

import (
	"html/template"
	"time"
)

type ExportParams struct {
	Dictionary     *Dictionary
	LanguageName   string
	Copyright      string
	AuthorsNote    string
	OutputPath     string
	UseEmbeddedCSS bool
	CSS            template.CSS
	CSSFile        string
	Multipage      bool
	ShowNavbar     bool
	IndexPage      bool
	HTMLEntries    []template.HTML
	NavbarHTML     template.HTML
	Timestamp      time.Time
	GenerationTime time.Duration
	NumWords       int
	Author         string
}

// Create a default ExportParams object.
func NewExportParams(dict *Dictionary) *ExportParams {
	return &ExportParams{
		Dictionary:     dict,
		LanguageName:   dict.LanguageName,
		Copyright:      "No copyright information provided.",
		UseEmbeddedCSS: true,
		CSS:            template.CSS(CSS), // Default CSS
		Timestamp:      time.Now(),
		NumWords:       len(dict.Entries),
	}
}

// Converts an ExportParams object into a map[string]any for passing them to
// HTML templates.
func (p *ExportParams) ToTemplateParams() map[string]any {
	return map[string]any{
		"LanguageName":   p.LanguageName,
		"HTMLEntries":    p.HTMLEntries,
		"Timestamp":      p.Timestamp,
		"NumWords":       p.NumWords,
		"GenerationTime": p.GenerationTime,
		"Copyright":      template.HTML(p.Copyright),
		"AuthorsNote":    template.HTML(p.AuthorsNote),
		"UseEmbeddedCSS": p.UseEmbeddedCSS,
		"CSS":            p.CSS,
		"CSSFile":        template.HTML(p.CSSFile),
		"Multipage":      p.Multipage,
		"ShowNavbar":     p.ShowNavbar,
		"NavbarHTML":     p.NavbarHTML,
		"IndexPage":      p.IndexPage,
		"Author":         p.Author,
	}
}
