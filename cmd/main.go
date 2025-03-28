package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage:   "A program for managing conlang lexicons",
		Authors: []*cli.Author{{Name: "Lemuria"}},
		Commands: []*cli.Command{
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "Import a lexicon from another format.",
				Action:  cmdImport,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Usage: "Format of the file to import from", Required: true, Aliases: []string{"f"}},
					&cli.StringFlag{Name: "input", Usage: "File to import from", Required: true, Aliases: []string{"i"}},
					&cli.StringFlag{Name: "output", Usage: "File to output LLEX json to. lp for Lexique Pro .db files, the only supported file format.", Aliases: []string{"o"}},
					&cli.StringFlag{Name: "language-name", Usage: "Name of the language to be imported (Lexique Pro does not include the name)"},
				},
			},
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "Export a lexicon.",
				Action:  cmdExport,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Usage: "Format to export into", Required: true, Aliases: []string{"f"}},
					&cli.StringFlag{Name: "input", Usage: "File to import from", Required: true, Aliases: []string{"i"}},

					// We call it the output path, because the format can either be a single file or a directory (in the case of a website export.)
					&cli.StringFlag{Name: "output", Usage: "Path to output the exported lexicon to.", Aliases: []string{"o"}},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
