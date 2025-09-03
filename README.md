# Skytale

*Skytale was a cipher tool used by the ancient Greeks and Spartans to encrypt and decrypt messages during military campaigns: (https://en.wikipedia.org/wiki/Scytale)*

## Description

The objective of this project is to create a simple command-line application that can encrypt and decrypt files via AES-GCM.

#### *Note on stages of files in program*:

unencrypted (original file) -> encrypted -> decrypted

While their is no technical difference between a *decrypted* file and an *unencrypted* file, for the purposes of documentation, it is useful to differentiate between an original *unencrypted* file and a *decrypted* file that has already stepped through the program.

### Requirements

- Must set an enviroment variable "AES_KEY" with the path of the file containing the 16 or 32 bytes used for encrypting/decrypting files. **Key file must have file extension *.aes*.**
- **Encrypted files to be decrypted must have file extension *.enc*** in order for decrypt option to recognize the file as valid candidate to be decrypted.

### Installation

```
# Copy the source code from github
$ git clone https://github.com/Ty-Grisham/skytale

# Navigate to the cmd directory of the repository
~/skytale$ cd cmd

# Build application with Go's compiler
~/skytale/cmd$ go build -o skytale

# Add the go install directory to your system's shell path
$ export PATH=$PATH:/path/to/your/install/directory
```

### Features

- Encrypts file given filepath to unencrypted file and encryption key
- Decrypts file given filepath to encrypted file and key used to encrypt file

### Usage
```
# Encrypt file
$ skytale -e </path/to/unencrypted/file.pdf>

# Decrypt file
$ shytale -d </path/to.encrypted/file.enc>
```

*Note: Both decrypt and encrypt options will create a new encrypted/decrypted file, but will not delete previous versions of the file.*

### Example

The example file to be encrypted is `sensitive_info.txt`. The file is currently in its *unencrypted* state because it has not been encrypted/decrypted yet.

Before the file can be encrypted, the environment variable "AES_KEY" must be set to the path of the 32 (or 16) byte encryption key. **The key must end with the file extension *.aes*.** This can be accomplished by executing the following command:

```
# Setting environment variable
$ export AES_KEY=<path/to/key.aes>
```

The following command encrypts `sensitive_info.txt`:

```
# Encrypting file
$ skytale -e <path/to/sensitive_info.txt>
```

Encrypting the file creates a new file, `sensitive_info.txt.enc`, in the same directory of the unencrypted file, `sensitive_info.txt`. The unencrypted file was not removed by the code, so it is the user's discretion whether to remove it or not.

Now, the file should be decrypted in order to read it. With the environment variable "AES_KEY" still set to the key that encrypted the file, decrypting the file can be accomplished with the following command:

```
# Decrypting file
$ skytale -d <path/to/sensitive_info.txt.enc>
```

Again, this will create a new decrypted file, `DECRYPTED-sensitive_info.txt`, in the same directory as the encrypted file. Notice that the file extension for encrypted files *.enc* was stripped away to leave the file extension of the original unencrypted file, in this case *.txt*. This was done so that the file may be viewed in whatever application the user prefers for viewing the file type.

Once finished viewing, the user should manually remove `DECRYPTED-sensitive_info.txt` so that there is no way of viewing the sensitive info in the file without decrypting `sensitive_info.txt.enc`.

```
# Removing decrypted file
/path/to$ rm DECRYPTED-sensitive_info.txt 
```
