package llex

import (
	"encoding/json"
	"os"
)

func ReadDictionary(path string) (*Dictionary, error) {
	jsonText, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dict Dictionary

	err = json.Unmarshal(jsonText, &dict)
	return &dict, err
}
