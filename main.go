package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	description = "A basic command-line XOR cipher tool."
)

var (
	key      string
	isHex    bool
	isBase64 bool
)

func usage() {
	fmt.Printf("%s: %s\n", os.Args[0], description)
	fmt.Printf("Usage: %s [-h|-b] key\n", os.Args[0])
	flag.PrintDefaults()
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
	os.Exit(-1)
}

func init() {
	flag.Usage = usage
	flag.BoolVar(&isHex, "h", false, "Specifies the key is in hex format.")
	flag.BoolVar(&isBase64, "b", false, "Specifies the key is base64 encoded.")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(-1)
	}

	keyStr := args[0]

	var key []byte
	var err error

	formatCount := 0
	formats := []bool{isHex, isBase64}
	for _, format := range formats {
		if format {
			formatCount++
		}
	}

	if formatCount > 1 {
		die(fmt.Errorf("only one key format can be specified."))
	}

	if isHex {
		key, err = hex.DecodeString(keyStr)
	} else if isBase64 {
		key, err = base64.StdEncoding.DecodeString(keyStr)
	} else {
		key = []byte(keyStr)
	}

	if err != nil {
		die(err)
	}

	if len(key) == 0 {
		die(fmt.Errorf("key length cannot be zero."))
	}

	buf := make([]byte, 4096)

	k := 0
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			die(err)
		}
		for i := 0; i < n; i++ {
			buf[i] ^= key[k%len(key)]
			k++
		}
		n, err = os.Stdout.Write(buf[:n])
		if err != nil {
			die(err)
		}
	}
}
