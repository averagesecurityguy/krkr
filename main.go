package main

import (
	"flag"
	"fmt"
	"os"
)

type cracker interface {
	Load(string) []string
	Hash(string, string)
}

func main() {
	var hashType string
	var hashFile string
	var wordFile string
	var c cracker

	flag.StringVar(&hashType, "t", "none", "The type of hash to crack.")
	flag.StringVar(&wordFile, "w", "words.txt", "Password list")
	flag.StringVar(&hashFile, "f", "hashes.txt", "File containing password hashes.")

	flag.Parse()

	// Set our loader and hasher functions based on the provided hash type
	switch hashType {
	case "ansible-vault":
		c = new(AnsibleVault)
	case "mongo-scram":
		c = new(MongoScram)
	case "mongo-cr":
		c = new(MongoCR)
	default:
		fmt.Println("Invalid hash type.")
		flag.Usage()
		os.Exit(0)
	}

	// Attempt to load our password hashes using the defined loader.
	hashList := c.Load(hashFile)
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

	fmt.Printf("Loaded %d hashes.\n", len(hashList))
	fmt.Printf("Loaded %d words.\n", len(wordList))

	// Attempt to crack our passwords. Kick off one Go routine for each candidate
	// password.
	var signal = make(chan bool)
	for i := range wordList {
		go func(word string) {
			for i := range hashList {
				c.Hash(hashList[i], word)
				signal <- true
			}
		}(wordList[i])
	}

	for i:=0; i < (len(wordList) * len(hashList)); i++ {
		<-signal
	}
}
