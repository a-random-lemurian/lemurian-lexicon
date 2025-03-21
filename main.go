package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	language, err := ImportFromLexiquePro("Kenahari.db")
	if err != nil {
		log.Fatal(err)
	}

	if v, err := json.Marshal(language); err == nil {
		fmt.Printf(string(v))
	}

	return
}
