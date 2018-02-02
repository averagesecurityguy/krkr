package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strings"
)

type Hash struct {
	User string
	Salt []byte
	Key  string
}

func parseMongoScramHash(hash string) string {
	parts := strings.Split(hash, ":")
	user := parts[0]

	parts = strings.Split(parts[2], "$")
	salt := parts[3]
	expected := parts[4]

	for {
		if len(salt)%3 == 0 {
			break
		}
		salt = salt + "="
	}

	//Decode our salt
	decoded, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%s:%s", user, decoded, expected)
}

// Parse the given file to get the mongo-scram hashes
func loadMongoScramHashes(filename string) []string {
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
			hash := parseMongoScramHash(text)
			hashes = append(hashes, hash)
		}
	}

	return hashes
}

func calculateMongoScramHash(hash, password string) {
	/*
	   Calculate the MongoDB SCRAM-SHA-1 hash. It varies from the standard
	   slightly by calculating the MD5 of the password and hex encoding it before
	   putting it through the PBKDF2 function.

	   Thanks @StrangeWill for helping me with that.
	*/
	params := strings.Split(hash, ":")
	user := params[0]
	salt := []byte(params[1])
	target := params[2]

	str := fmt.Sprintf("%s:mongo:%s", user, password)
	pwd_md5 := hex.EncodeToString(md5.New().Sum([]byte(str)))
	salted_password := pbkdf2.Key([]byte(pwd_md5), salt, 10000, 20, sha1.New)

	client_key := hmac.New(sha1.New, salted_password)
	client_key.Write([]byte("Client Key"))

	stored_key := sha1.New()
	stored_key.Write(client_key.Sum(nil))

	calculated := base64.StdEncoding.EncodeToString(stored_key.Sum(nil))

	if calculated == target {
		fmt.Printf("%s:%s\n", target, password)
	}
}
