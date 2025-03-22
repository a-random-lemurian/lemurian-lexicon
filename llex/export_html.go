package llex

import (
	"bytes"
	"html/template"
	"log"
	"sort"
)

var HtmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>Dictionary</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>
    <style>
        .dictionary ol {
            margin: 0px;
        }
    </style>
</head>
<body style="max-width: 600px; margin: auto;">
    <h1>{{.LanguageName}}</h1>
    <hr>
    <div class="dictionary">
	{{range .Entries}}
	{{if .Word}}
	<div class="entry">
	<b class="headword">{{.Word}}</b> <i class="part-of-speech">{{.POS}}</i> <br>
	<ol>
		{{range .Definitions}}
		<li>{{.Text}}</li>
		{{end}}
	</ol>
	</div>
	{{else}}
	{{end}}
	{{end}}
    </div>
    <hr>
</body>
</html>
`

type htmlParameters struct {
	LanguageName string
	Entries      []*Entry
}

// Export a Dictionary to a single HTML file.
func ExportSinglePageHTML(dict *Dictionary) (string, error) {
	var html bytes.Buffer

	sortedEntries := dict.Entries
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].Word < sortedEntries[j].Word
	})

	t, err := template.New("html").Parse(HtmlTemplate)
	if err != nil {
		log.Print(err)
		return "", err
	}

	if err := t.Execute(&html, htmlParameters{
		LanguageName: dict.LanguageName,
		Entries:      sortedEntries,
	}); err != nil {
		log.Print(err)
		return "", err
	}

	return html.String(), nil
}

/*
func ExportStaticHTML(dict *Dictionary) (string, error) {

}
*/
