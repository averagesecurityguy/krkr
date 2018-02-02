# cracker
A Go based password cracker for odd hashes.

## Building

* `git clone https://github.com/averagesecurityguy/cracker`
* `cd cracker`
* `go build`

## Usage

```
Usage of ./cracker:
  -f string
    	File containing password hashes. (default "hashes.txt")
  -t string
    	The type of hash to crack. (default "none")
  -w string
    	Password list (default "words.txt")
```

# Supported Hashes

 * Mongodb SCRAM-SHA1
 * Mongodb CR
 * Ansible Vault
