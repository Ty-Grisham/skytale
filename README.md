# Skytale

The objective of this project is to create a simple command-line application that can encrypt and decrypt files.

#### *Note on stages of files in program*:

unencrypted (original file) -> encrypted -> decrypted

While their is no technical difference between a *decrypted* file and an *unencrypted* file, for the purposes of documentation, it is useful to differentiate between an original *unencrypted* file and a *decrypted* file that has already stepped through the program.

## Usage

This program will be used to encrypt and decrypt files given an encryption key (the key must have a *.aes* file extension).

One encrypts a file executing the following command in a terminal (linux): 
```
$PATH skytale -e <sensitive_info.txt>
```