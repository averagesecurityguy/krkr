package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

const threads = 16

type candidate struct {
	word string
	hash string
}

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

	// Load our password candidates into a channel
	var candidates = make(chan *candidate, len(wordList) * len(hashList))
	for _, w := range wordList {
		for _, h := range hashList {
			candidates <- &candidate{word: w, hash: h}
		}
	}
	close(candidates)

	fmt.Printf("Loaded %d hashes.\n", len(hashList))
	fmt.Printf("Loaded %d words.\n", len(wordList))

	// Listen for Ctrl-C and kill the program.
	var sig = make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt)
    go func() {
        for _ = range sig {
			os.Exit(1)
		}
    }()

	// Start our threads for processing hashes.
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for cnd := range candidates {
				c.Hash(cnd.hash, cnd.word)
			}
		}()
	}

	// Wait for our workers to finish.
	wg.Wait()
}
