package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
)

func strSliceCommaList(slice []string) string {
	output := ""
	for i, item := range slice {
		output += item
		if i < len(slice)-1 {
			output += ", "
		}
	}
	return output
}

func cmdListFormats(cCtx *cli.Context) error {
	var supportedFormatString string

	if cCtx.Bool("json") {
		type formatObject struct {
			Export []string `json:"export"`
			Import []string `json:"import"`
		}
		supportedFormatsJSON, err := json.Marshal(formatObject{Export: supportedExportFormats, Import: supportedImportFormats})
		if err != nil {
			return err
		}
		supportedFormatString = string(supportedFormatsJSON)
	} else {
		supportedFormatString += "Supported import formats: " + strSliceCommaList(supportedImportFormats) + "\n"
		supportedFormatString += "Supported export formats: " + strSliceCommaList(supportedExportFormats)
	}

	fmt.Println(supportedFormatString)

	return nil
}
