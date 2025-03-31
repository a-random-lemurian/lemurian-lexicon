package llex

import (
	"bytes"
	"html/template"
	"maps"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

// Parameters for how a list of entries should be split based on their first letter.
//
// TODO: Allow the user to set case sensitivity instead of hardcoding it to false.
type splitWordParams struct {
	Entries       []*Entry // The entries to be split.
	CaseSensitive bool     // Whether words with different casings are treated differently.
}

// Split a dictionary into a map of Entry slices based on the word's first letter.
//
// TODOs that we are not going to deal with right now because Kenahari doesn't need
// them:
//
//   - The ability to split by the first two letters, for when lexicons start to get
//     really long
//   - Actual unicode support - when testing two letter splits, a question mark character
//     appeared in a word file containing one word that started with ká (note the diacritic).
func splitWordsByLetter(params *splitWordParams) map[string][]*Entry {
	entries := params.Entries
	alphabeticalMap := make(map[string][]*Entry)

	for _, entry := range entries {
		// Ideally, empty entries (entry objects where the word is a blank string)
		// would be caught and removed at earlier stages.
		if entry.Word == "" {
			continue
		}

		firstLetter := entry.Word[0:1]
		if !params.CaseSensitive {
			firstLetter = strings.ToLower(firstLetter)
		}

		alphabeticalMap[firstLetter] = append(alphabeticalMap[firstLetter], entry)
	}

	return alphabeticalMap
}

// Convenience function to create a file with a string.
func writeStringToFile(data string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(data)
	return nil
}

var navbarTemplate = `<nav class="navbar">
{{range .Letters}}<a class="letter" id="nav-letter-{{.}}" href="./{{.}}.html">{{.}}</a>
{{end}}</nav>`

func generateNavbarHtml(words *map[string][]template.HTML) (template.HTML, error) {
	t, err := template.New("navbar").Parse(navbarTemplate)
	if err != nil {
		return "", err
	}

	type navbarHtmlParameters struct {
		Letters []string
	}
	var params navbarHtmlParameters
	for letter := range maps.Keys(*words) {
		params.Letters = append(params.Letters, letter)
	}
	slices.Sort(params.Letters)

	var html bytes.Buffer
	err = t.Execute(&html, params)
	if err != nil {
		return "", err
	}

	return template.HTML(html.String()), nil
}

// Export a Dictionary to a static set of HTML files.
func ExportStaticHTML(params *ExportParams) error {
	startTime := time.Now()
	outdir := params.OutputPath
	CSS_FILE := "index.css"

	// Create the output directory.
	err := os.MkdirAll(outdir, 0755)
	if err != nil {
		return err
	}

	// Split the entry list into individual lists by the first letter.
	alphabeticalMap := splitWordsByLetter(&splitWordParams{
		Entries:       params.Dictionary.Entries,
		CaseSensitive: false,
	})

	// Now, prepare the HTML strings for writing.
	alphabeticalMapHTML := make(map[string][]template.HTML)
	for letter, entries := range alphabeticalMap {
		sortEntries(entries)
		entryHTMLSlice, err := batchGenerateEntryHTML(entries)
		if err != nil {
			return err
		}
		alphabeticalMapHTML[letter] = entryHTMLSlice
	}

	// Generate the navigation bar, which will allow users to navigate
	// by letter.
	navbarHTML, err := generateNavbarHtml(&alphabeticalMapHTML)
	if err != nil {
		return err
	}

	params.NavbarHTML = navbarHTML
	params.ShowNavbar = true
	params.Multipage = true
	params.GenerationTime = time.Since(startTime)
	params.IndexPage = true

	// Generate index.html.
	indexHTML, err := executeHTMLTemplate(params.ToTemplateParams())
	if err != nil {
		return err
	}

	params.IndexPage = false

	err = writeStringToFile(indexHTML, path.Join(outdir, "index.html"))
	if err != nil {
		return err
	}

	// Begin generating the HTML pages.
	for letter, entries := range alphabeticalMapHTML {
		params.HTMLEntries = entries
		params.NumWords = len(entries)
		htmlString, err := executeHTMLTemplate(params.ToTemplateParams())
		if err != nil {
			return err
		}

		err = writeStringToFile(htmlString, path.Join(outdir, letter+".html"))
		if err != nil {
			return err
		}
	}

	// Write the CSS file out.
	err = writeStringToFile(CSS, path.Join(outdir, CSS_FILE))
	if err != nil {
		return err
	}

	params.CSSFile = CSS_FILE
	params.Multipage = true

	// Generate the all-words.html file.
	allHTML, err := ExportSinglePageHTML(params)
	if err != nil {
		return err
	}

	return writeStringToFile(allHTML, path.Join(outdir, "all-words.html"))
}
