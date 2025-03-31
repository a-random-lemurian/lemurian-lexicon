package llex

type QualifiedStrings struct {
	Qualifiers []string `json:"qualifiers,omitempty"`
	Text       string   `json:"text"`
}

type Definition QualifiedStrings
type IPA QualifiedStrings

type Entry struct {
	Word           string        `json:"word"`
	POS            string        `json:"partOfSpeech"`
	Pronunciations []*IPA        `json:"pronunciations,omitempty"`
	Definitions    []*Definition `json:"definitions"`
	UsageNotes     []string      `json:"usageNotes,omitempty"`
	Etymology      string        `json:"etymology,omitempty"`
	BorrowedWord   string        `json:"borrowedWord,omitempty"`
}

type Dictionary struct {
	LanguageName string   `json:"languageName"`
	Entries      []*Entry `json:"entries"`
}

type StaticExportParams struct {
	Dictionary  *Dictionary
	Copyright   string
	AuthorsNote string
}
