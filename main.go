package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var hashType string
	var hashFile string
	var wordFile string
	var loader func(filename string) []string
	var hasher func(hash, password string)

	flag.StringVar(&hashType, "t", "none", "The type of hash to crack.")
	flag.StringVar(&wordFile, "w", "words.txt", "Password list")
	flag.StringVar(&hashFile, "f", "hashes.txt", "File containing password hashes.")

	flag.Parse()

	// Set our loader and hasher functions based on the provided hash type
	switch hashType {
	case "ansible-vault":
		loader = loadAnsibleVaultHashes
		hasher = calculateAnsibleVaultHash
	case "mongo-scram":
		loader = loadMongoScramHashes
		hasher = calculateMongoScramHash
	case "mongo-cr":
		loader = loadMongoCrHashes
		hasher = calculateMongoCrHash
	default:
		fmt.Println("Invalid hash type.")
		flag.Usage()
		os.Exit(0)
	}

	// Attempt to load our password hashes using the defined loader.
	hashList := loader(hashFile)
	if len(hashList) == 0 {
		fmt.Println("No password hashes loaded.")
		os.Exit(0)
	}

	// Attempt to load our password list.
	wordList := loadWords(wordFile)
	if len(wordList) == 0 {
		fmt.Println("No passwords loaded.")
		os.Exit(0)
	}

	fmt.Printf("Loaded %d words.\n", len(wordList))

	// Attempt to crack our passwords. Kick off one Go routine for each candidate
	// password.
	var signal = make(chan bool)
	for i := range wordList {
		go func(word string) {
			for i := range hashList {
				hasher(hashList[i], word)
				signal <- true
			}
		}(wordList[i])
	}

	for range wordList {
		<-signal
	}
}
