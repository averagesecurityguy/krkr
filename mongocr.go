package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func loadMongoCrHashes(filename string) []string {
	var hashes []string

	data := readFile(filename)

	for _, line := range strings.Split(data, "\n") {
		if line != "" {
			hashes = append(hashes, line)
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
