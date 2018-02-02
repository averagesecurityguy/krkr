package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func loadMongoCrHashes(filename string) []string {
	var hashes []string

	data, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not load hashes: %s\n", filename)
	}

	defer data.Close()

	scan := bufio.NewScanner(data)
	for scan.Scan() {
		text := scan.Text()
		if text != "" {
			hashes = append(hashes, text)
		}
	}

	return hashes
}

func calculateMongoCrHash(hash, password string) {
	parts := strings.Split(hash, ":")
	user := parts[0]
	target := parts[1]

	str := fmt.Sprintf("%s:mongo:%s", user, password)
	pwd_md5 := md5.New().Sum([]byte(str))

	calculated := hex.EncodeToString(pwd_md5)

	if target == calculated {
		fmt.Printf("%s:%s\n", hash, password)
	}

}
