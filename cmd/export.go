package main

import (
	"bufio"
	"encoding/json"
	"html"
	"os"
	"slices"

	"github.com/a-random-lemurian/lemurian-lexicon/llex"
	"github.com/urfave/cli/v2"
)

var supportedExportFormats = []string{
	"html",    // Single-file HTML
	"website", // Static website or directory
}

// Get the contents of a file that can be passed in through command-line arguments.
// Used for optional files.
//
// If filename is "", an empty string and a nil error will be returned.
func getFileOptional(filename string) ([]byte, error) {
	if filename == "" {
		return []byte{}, nil
	}
	return os.ReadFile(filename)
}

func escapeHtml(text string, escape bool) string {
	if escape {
		text = html.EscapeString(text)
	}
	return text
}

func cmdExport(cCtx *cli.Context) error {
	exportFmt := cCtx.String("format")
	inputFile := cCtx.String("input")
	outputPath := cCtx.String("output")

	if !slices.Contains(supportedExportFormats, exportFmt) {
		return &ErrorUnsupportedFormat{attemptedFormat: exportFmt}
	}

	dictionaryRawJson, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Attempt to create the output file before starting the generation
	// process, so that if there is a problem with the output file, time
	// is not wasted generating a result that will never be written.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	var dictionary llex.Dictionary
	err = json.Unmarshal(dictionaryRawJson, &dictionary)
	if err != nil {
		return err
	}

	var output string

	params := llex.NewStaticExportParams(&dictionary)

	// Retrieve authors' note and copyright text.

	treatAsHtml := cCtx.Bool("treat-as-html")
	
	// Only overwrite params.Copyright so that if the copyright file is an empty string,
	// the default single-file HTML export will still explicitly state that no copyright
	// information was provided.
	copyrightBytes, err := getFileOptional(cCtx.String("copyright"))
	if err != nil {
		return err
	}
	copyright := string(copyrightBytes)
	copyright = escapeHtml(copyright, treatAsHtml)
	if copyright != "" {
		params.Copyright = copyright
	}

	authorsNoteBytes, err := getFileOptional(cCtx.String("authors-note"))
	authorsNote := escapeHtml(string(authorsNoteBytes), treatAsHtml)
	if err != nil {
		return err
	}
	params.AuthorsNote = authorsNote

	switch exportFmt {
	case "html":
		output, err = llex.ExportSinglePageHTML(params)
	}

	if err != nil {
		return err
	}
	
	_, err = writer.WriteString(output)
	if err != nil {
		return err
	}

	return writer.Flush()
}
