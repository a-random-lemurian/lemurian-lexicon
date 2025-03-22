package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "Import a lexicon from another format.",
				Action:  cmdImport,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Usage: "Format of the file to import from", Required: true, Aliases: []string{"t"}},
					&cli.StringFlag{Name: "file", Usage: "File to import from", Required: true, Aliases: []string{"f"}},
					&cli.StringFlag{Name: "output", Usage: "File to output LLEX json to. lp for Lexique Pro .db files, the only supported file format.", Aliases: []string{"o"}},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
