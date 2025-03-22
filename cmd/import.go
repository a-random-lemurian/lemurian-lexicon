package main

import (
	"os"
	"fmt"
	"bufio"
	
	"encoding/json"

	"github.com/a-random-lemurian/lemurian-lexicon/llex"
	"github.com/urfave/cli/v2"
)

type ErrorUnsupportedFormat struct {
	attemptedFormat string
}

func (e *ErrorUnsupportedFormat) Error() string {
	return "unsupported import format '" + e.attemptedFormat + "'"
}

func cmdImport(cCtx *cli.Context) error {
	importFmt := cCtx.String("format")
	importFile := cCtx.String("file")
	outputFile := cCtx.String("output")

	// Todo: Do not hardcode once extra formats added
	if importFmt != "lp" {
		return &ErrorUnsupportedFormat{attemptedFormat: importFmt}
	}

	dict, err := llex.ImportFromLexiquePro(importFile)
	if err != nil {
		return err
	}
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

	writer := bufio.NewWriter(file)
	_, err = writer.Write(dictJson)
	return err
}
