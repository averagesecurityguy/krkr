package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func parseAnsibleVaultHash(data string) string {
	lines := strings.Split(data, "\n")
	hexData := strings.Join(lines[1:], "")
	decoded := decodeHex(hexData)
	split := strings.Split(string(decoded), "\n")

	return fmt.Sprintf("%s:%s:%s", split[0], split[1], split[2])
}

func loadAnsibleVaultHashes(filename string) []string {
	var hashes []string

	data := readFile(filename)

	if !strings.HasPrefix(data, "$ANSIBLE_VAULT;1.1;AES256") {
		fmt.Println(data)
		fmt.Println("File contains an invalid hash.")
		return hashes
	}

	vh := parseAnsibleVaultHash(data)
	hashes = append(hashes, vh)

	return hashes
}

func calculateAnsibleVaultHash(hash, password string) {
	params := strings.Split(hash, ":")
	salt := decodeHex(params[0])
	target := params[1]
	data := decodeHex(params[2])

	dk := pbkdf2.Key([]byte(password), salt, 10000, 80, sha256.New)

	mac := hmac.New(sha256.New, dk[32:64])
	mac.Write(data)
	sum := mac.Sum(nil)
	candidate := hex.EncodeToString(sum)

	if target == candidate {
		fmt.Printf("%s:%s\n", hash, password)
	}
}
