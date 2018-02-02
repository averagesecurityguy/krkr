package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
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

	wordlist, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open wordlist: %s\n", filename)
		return words
	}

	defer wordlist.Close()

	scan := bufio.NewScanner(wordlist)
	for scan.Scan() {
		text := scan.Text()

		if strings.HasPrefix(text, "#") {
			continue
		}

		words = append(words, text)
	}

	return words
}

func testWords(words <-chan string, hashes []string, signal chan<- bool) {

}
