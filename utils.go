package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
)

func decodeHex(hexData string) []byte {
	decoded, err := hex.DecodeString(hexData)
	if err != nil {
		fmt.Println("Decode Hex: %s", err)
	}

	return decoded
}

func loadWords(filename string) []string {
	var words []string

	data := readFile(filename)
	for _, line := range strings.Split(data, "\n") {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		words = append(words, line)
	}

	return words
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not read file: %s\n", filename)
		return ""
	}

	data := string(bytes)

	if len(data) == 0 {
		fmt.Println("File contains no data.")
		return ""
	}

	return data
}
