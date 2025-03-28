package llex

import (
	"bufio"
	"os"
	"strings"
)

// Import definitions from a Lexique Pro file.
func ImportFromLexiquePro(filename string) (*Dictionary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dictionary := &Dictionary{}

	scanner := bufio.NewScanner(file)

	var currentEntry *Entry

	for scanner.Scan() {
		text := scanner.Text()
		tokens := strings.Split(text, " ")

		switch tokens[0] {
		case `\lx`:
			if currentEntry != nil {
				dictionary.Entries = append(dictionary.Entries, currentEntry)
			}

			currentEntry = &Entry{Word: strings.Join(tokens[1:], " ")}
			currentEntry.Definitions = make([]*Definition, 0)
			currentEntry.Pronunciations = make([]*IPA, 0)
			currentEntry.UsageNotes = make([]string, 0)
		case `\ps`:
			currentEntry.POS = strings.Join(tokens[1:], " ")
		case `\de`, `\ge`:
			definitions := strings.Split(strings.Join(tokens[1:], " "), ";")
			for _, def := range definitions {
				currentEntry.Definitions = append(currentEntry.Definitions, &Definition{
					Text: def,
				})
			}
		case `\et`:
			currentEntry.Etymology = strings.Join(tokens[1:], " ")
		default:
			continue
		}

	}
	
	// Add the final entry to the list.
	dictionary.Entries = append(dictionary.Entries, currentEntry)

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dictionary, nil
}
