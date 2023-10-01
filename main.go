package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

const (
	description = `A basic command-line XOR cipher tool.`
	bufSize     = 1024 * 16
)

var (
	key          []byte
	keyInput     string
	isFile       bool
	isHex        bool
	isBase64     bool
	isDecrypting bool
)

func usage() {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], description)
	fmt.Fprintf(os.Stderr, "Usage: %s [options] -k key [files...]\n", os.Args[0])
	flag.PrintDefaults()
}

func handlePanic() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(-1)
	}
}

func init() {
	flag.Usage = usage
	flag.StringVarP(&keyInput, "key", "k", "", "The input key/file. (required)")
	flag.BoolVarP(&isBase64, "base64", "b", false, "Specifies the key is base64 encoded.")
	flag.BoolVarP(&isHex, "hex", "h", false, "Specifies the key is in hex format.")
	flag.BoolVarP(&isFile, "file", "f", false, "Specifies the input key is a file.")
	flag.BoolVarP(&isDecrypting, "decrypt", "d", false, "Decrypt files with a .xor extension.")
	flag.Parse()
}

func main() {
	defer handlePanic()

	if keyInput == "" {
		flag.Usage()
		os.Exit(1)
	}

	files := flag.Args()

	var err error

	formatCount := 0
	formats := []bool{isHex, isBase64, isFile}
	for _, format := range formats {
		if format {
			formatCount++
		}
	}

	if formatCount > 1 {
		panic("only one key format can be specified.")
	}

	if isFile {
		f, err := os.Open(keyInput)
		if err != nil {
			panic(fmt.Errorf("failed to open key file: %w", err))
		}
		defer f.Close()
		key, err = io.ReadAll(f)
		if err != nil {
			panic(fmt.Errorf("failed to read key file: %w", err))
		}
	} else if isHex {
		key, err = hex.DecodeString(keyInput)
	} else if isBase64 {
		key, err = base64.StdEncoding.DecodeString(keyInput)
	} else {
		key = []byte(keyInput)
	}

	if err != nil {
		panic(err)
	}

	if len(key) == 0 {
		panic("key length cannot be zero.")
	}

	if len(files) > 0 {
		for _, file := range files {
			err := xorFile(file)
			if err != nil {
				panic(err)
			}
		}
	} else {
		xorStdin()
	}
}

func xorStdin() {
	buf := make([]byte, bufSize)

	k := 0
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		for i := 0; i < n; i++ {
			buf[i] ^= key[k]
			k = (k + 1) % len(key)
		}
		n, err = os.Stdout.Write(buf[:n])
		if err != nil {
			panic(err)
		}
	}
}

func xorFile(inPath string) error {
	fmt.Fprintf(os.Stderr, "%s: ", inPath)

	info, err := os.Stat(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed\n")
		return err
	}
	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "skipping directory\n")
		return nil
	}

	isEncrypted := strings.HasSuffix(strings.ToLower(inPath), ".xor")
	if isEncrypted != isDecrypting {
		fmt.Fprintf(os.Stderr, "skipping\n")
		return nil
	}

	outPath := inPath
	if isEncrypted {
		outPath = outPath[:len(inPath)-4]
	} else {
		outPath += ".xor"
	}

	if _, err := os.Stat(outPath); err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "failed\n")
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "output file exists\n")
		return nil
	}

	fin, err := os.Open(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open input file\n")
		return err
	}
	defer fin.Close()

	fout, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR, info.Mode().Perm())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open output file\n")
		return err
	}
	defer fout.Close()

	buf := make([]byte, bufSize)
	k := 0
	for {
		n, err := fin.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "failed to read input file\n")
			return err
		}
		if n == 0 {
			break
		}
		for i := 0; i < n; i++ {
			buf[i] ^= key[k]
			k = (k + 1) % len(key)
		}
		fout.Write(buf[:n])
	}

	fmt.Fprintf(os.Stderr, "%s\n", outPath)
	return nil
}
