# xor

A basic command-line XOR cipher tool.

# Usage

`xor [options] -k key [files...]`

- `-k` `--key` specifies the input key / key file.
- Only one of the following key formats can be specified:
  - `-h` `--hex` specifies the key is in hex format.
  - `-b` `--base64` specifies the key is base64 encoded.
  - `-f` `--file` specifies the key is a file.
- `-d` `--decrypt` decrypts input files with a `.xor` extension. Has no effect when processing from standard input.

## Process from standard input
```sh
$ echo Plaintext > plaintxt
$ cat plaintxt | xor -k cipherkey > ciphertxt
$ xxd ciphertxt
00000000: 3305 1101 0b06 0e1d 0d69                 3........i
$ cat ciphertxt | xor -k cipherkey
Plaintext

# base64 key
$ echo -n cipherkey | base64
Y2lwaGVya2V5
$ cat ciphertxt | xor -bk Y2lwaGVya2V5
Plaintext

# hex key
$ echo -n cipherkey | xxd -p
6369706865726b6579
$ cat ciphertxt | xor -hk 6369706865726b6579
Plaintext

# key file
$ echo -n cipherkey > cipher.key
$ cat ciphertxt | xor -fk cipher.key
Plaintext
```

## Process files
Encrypted files will be saved with a `.xor` suffix.\
When encrypting, files with a `.xor` suffix will be ignored.\
Likewise, files without a `.xor` suffix will be ignored when decrypting.
```sh
$ for x in {a..c}; do echo "Plaintext $x" > "$x.txt"; done
$ ls
a.txt  b.txt  c.txt
$ xor -k cipherkey *.txt
a.txt: a.txt.xor
b.txt: b.txt.xor
c.txt: c.txt.xor
$ rm *.txt
$ ls
a.txt.xor  b.txt.xor  c.txt.xor
$ for f in *.xor; do xxd "$f"; done
00000000: 3305 1101 0b06 0e1d 0d43 087a            3........C.z
00000000: 3305 1101 0b06 0e1d 0d43 0b7a            3........C.z
00000000: 3305 1101 0b06 0e1d 0d43 0a7a            3........C.z
$ xor -d -k cipherkey *.xor
a.txt.xor: a.txt
b.txt.xor: b.txt
c.txt.xor: c.txt
$ rm *.xor
$ cat *.txt
Plaintext a
Plaintext b
Plaintext c
```

# Installation

Requires Go 1.21.1.
```sh
go install github.com/x65534/xor@latest
```
