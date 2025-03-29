package main

import (
	"bufio"
	"fmt"
	"os"

	"encoding/json"

	"github.com/a-random-lemurian/lemurian-lexicon/llex"
	"github.com/urfave/cli/v2"
)

var supportedImportFormats = []string{
	"lp",    // Lexique Pro database files
}

type ErrorUnsupportedFormat struct {
	attemptedFormat string
}

func (e *ErrorUnsupportedFormat) Error() string {
	return "unsupported format '" + e.attemptedFormat + "'"
}

func cmdImport(cCtx *cli.Context) error {
	importFmt := cCtx.String("format")
	importFile := cCtx.String("input")
	outputFile := cCtx.String("output")

	// Todo: Do not hardcode once extra formats added
	if importFmt != "lp" {
		return &ErrorUnsupportedFormat{attemptedFormat: importFmt}
	}

	dict, err := llex.ImportFromLexiquePro(importFile)
	if err != nil {
		return err
	}

	dict.LanguageName = cCtx.String("language-name")

	dictJson, err := json.Marshal(dict)
	if err != nil {
		return err
	}

	if outputFile == "" {
		fmt.Println(string(dictJson))
		return nil
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(dictJson)
	if err != nil {
		return err
	}

	return writer.Flush()
}
