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

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		writer := bufio.NewWriter(outputFile)
		_, err = writer.WriteString(html)
		if err != nil {
			return err
		}
	}

	return nil
}
