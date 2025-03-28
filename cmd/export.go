package main

import (
	"bufio"
	"encoding/json"
	"os"
	"slices"

	"github.com/a-random-lemurian/lemurian-lexicon/llex"
	"github.com/urfave/cli/v2"
)

var supportedExportFormats = []string{
	"html",    // Single-file HTML
	"website", // Static website or directory
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

	switch exportFmt {
	case "html":
		html, err := llex.ExportSinglePageHTML(&dictionary)
		if err != nil {
			return err
		}

		_, err = writer.WriteString(html)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
