package llex

import (
	"bytes"
	"html/template"
	"log"
	"sort"
	"time"
)

var HtmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>{{.LanguageName}} Dictionary</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>
    <style>
        .dictionary ol {
            margin: 0px;
        }
		a {
			color: #acc0fb;
		}
		body {
			background-color: #000;
			color: #fff;
			max-width: 600px;
			margin: auto;
		}
		.entry {
			break-inside: avoid;
			margin-bottom: 8px;
			padding-bottom: 8px;
		}
		.auxilliary p {
			font-size: 85%;
			margin: 0.1%;
		}
    </style>
</head>
<body>
    <h1><span class="language-name">{{.LanguageName}}</span></h1>
    <hr>
	{{if .AuthorsNote}}<div class="authors-note">{{.AuthorsNote}}</div>{{end}}
    <div class="dictionary">
	{{range .HTMLEntries}}
	{{.}}
	{{end}}
    </div>
    <hr>
	<p><b>Copyright</b>: {{.Copyright}}</p>
	<p>Generated by the
	   <a href="https://github.com/a-random-lemurian/lemurian-lexicon">Lemurian Lexicon Manager</a>
	   at <span class="timestamp">{{.Timestamp}}</span>.
	   Generation time <span class="generation-time">{{.GenerationTime}}</span>.
	   Contains <span class="num-words">{{.NumWords}}</span> words.</p>
</body>
</html>
`

var WordTemplate = `<div class="entry">
<b><span class="headword">{{.Word}}</span></b> <i><span class="part-of-speech">{{.POS}}</span></i> <br>
<ol class="definitions">
{{range .Definitions}}<li class="definition">{{.Text}}</li>{{end}}
</ol><div class="auxilliary">{{if .Etymology}}
<p>Etymology: <span class="etymology">{{.Etymology}}</span></p>{{else}}{{end}}
{{if .BorrowedWord}}<p>From: <span class="borrowed-from">{{.BorrowedWord}}</span>{{else}}{{end}}</div>
</div>`

type htmlParameters struct {
	LanguageName   string
	HTMLEntries    []template.HTML
	Timestamp      time.Time
	NumWords       int
	GenerationTime string
	Copyright 	   template.HTML
	AuthorsNote    template.HTML
}

type StaticExportParams struct {
	Dictionary  *Dictionary
	Copyright   string
	AuthorsNote string
}

func NewStaticExportParams(dict *Dictionary) *StaticExportParams {
	params := &StaticExportParams{
		Dictionary: dict,
	}
	params.Copyright = "No copyright information provided."
	params.AuthorsNote = ""
	return params
}

func (e *Entry) GenerateHTML() (string, error) {
	t, err := template.New("html").Parse(WordTemplate)
	if err != nil {
		return "", err
	}

	var html bytes.Buffer

	err = t.Execute(&html, e)

	if err != nil {
		return "", err
	}

	return html.String(), nil
}

func batchGenerateEntryHTML(entries []*Entry) ([]template.HTML, error) {
	var entriesHTML []template.HTML
	for _, entry := range entries {
		entryHTML, err := entry.GenerateHTML()
		if err != nil {
			return nil, err
		}
		entriesHTML = append(entriesHTML, template.HTML(entryHTML))
	}
	return entriesHTML, nil
}

func sortEntries(entries []*Entry) []*Entry {
	var sortedEntries []*Entry
	copy(entries, sortedEntries)
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].Word < sortedEntries[j].Word
	})
	return sortedEntries
}

// Export a Dictionary to a single HTML file.
func ExportSinglePageHTML(params *StaticExportParams) (string, error) {
	var html bytes.Buffer

	startTime := time.Now()
	dict := params.Dictionary

	sortedEntries := sortEntries(dict.Entries)
	sortedEntriesHTML, err := batchGenerateEntryHTML(sortedEntries)
		if err != nil {
			return "", err
	}

	t, err := template.New("html").Parse(HtmlTemplate)
	if err != nil {
		log.Print(err)
		return "", err
	}

	endTime := time.Now()

	if err := t.Execute(&html, htmlParameters{
		LanguageName:   dict.LanguageName,
		HTMLEntries:    sortedEntriesHTML,
		Timestamp:      endTime,
		GenerationTime: endTime.Sub(startTime).String(),
		NumWords:       len(dict.Entries),
		Copyright:      template.HTML(params.Copyright),
		AuthorsNote:    template.HTML(params.AuthorsNote),
	}); err != nil {
		log.Print(err)
		return "", err
	}

	return html.String(), nil
}

// Export a Dictionary to a static set of HTML files.
func ExportStaticHTML(params *StaticExportParams) {

}
