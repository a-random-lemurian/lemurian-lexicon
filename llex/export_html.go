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
    <div class="dictionary">
	{{range .HTMLEntries}}
	{{.}}
	{{end}}
    </div>
    <hr>
	<p>Generated by the Lemurian Lexicon Manager at <span class="timestamp">{{.Timestamp}}</span>.
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
<p>Etymology: <span class="etymology">{{.Etymology}}</span></p>{{else}}{{end}}</div>
</div>`

type htmlParameters struct {
	LanguageName   string
	HTMLEntries    []template.HTML
	Timestamp      time.Time
	NumWords       int
	GenerationTime string
}

type StaticExportParams struct {
	Dictionary *Dictionary
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

// Export a Dictionary to a single HTML file.
func ExportSinglePageHTML(params *StaticExportParams) (string, error) {
	var html bytes.Buffer

	startTime := time.Now()
	dict := params.Dictionary

	sortedEntries := dict.Entries
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].Word < sortedEntries[j].Word
	})

	var sortedEntriesHTML []template.HTML
	for _, entry := range sortedEntries {
		entryHTML, err := entry.GenerateHTML()
		if err != nil {
			return "", err
		}
		sortedEntriesHTML = append(sortedEntriesHTML, template.HTML(entryHTML))
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
	}); err != nil {
		log.Print(err)
		return "", err
	}

	return html.String(), nil
}

// Export a Dictionary to a static set of HTML files.
func ExportStaticHTML(params *StaticExportParams) {

}
