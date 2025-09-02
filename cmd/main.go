package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"

	cryptography "github.com/Ty-Grisham/skytale"
)

var (
	ErrInvalidExtension = errors.New("invalid file extension")
	ErrEnvNotSet        = errors.New("environment variable for encryption key not set")
	ErrEnvDoesNotExist  = errors.New("environment variable does not exist")
)

const (
	envVar = "AES_KEY" // name of environment variable for encryption key
	eExt   = ".enc"    // encyrption file extension
	keyExt = ".aes"    // key file extension
)

func main() {
	// Add command line options
	encrypt := flag.String("e", "", "File to be encrypted")
	decrypt := flag.String("d", "", "File to be decrypted")

	flag.Parse()

	// Acquire key from environment variable
	key, err := readKey(envVar, keyExt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	switch {
	// Encrypt file
	case *encrypt != "":
		cFilename, err := createEFile(*encrypt, eExt, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Created encrypted file: %s\n", cFilename)

	// Decrypt file
	case *decrypt != "":
		cFilename, err := createDFile(*decrypt, eExt, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Created decrypted file: %s\n", cFilename)

	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "invalid option")
		flag.Usage()
		os.Exit(1)
	}
}

// createEFile extracts the content from the input unencrypted filepath (uPath),
// and creates a new encrypted file with the given key and returns the
// path of the created encrypted file (string) a potential error. eExt represents
// the user-defined file extension of encrypted files.
func createEFile(uPath, eExt string, key []byte) (string, error) {
	// Extract unencrypted data from given file
	uData, err := os.ReadFile(uPath)
	if err != nil {
		return "", err
	}

	// Encrypt unencrypted data
	eData, err := cryptography.Encrypt(uData, key)
	if err != nil {
		return "", err
	}

	ePath := genEPath(uPath, eExt)

	// Write new encrypted file
	return ePath, os.WriteFile(ePath, eData, 0644)
}

// createDFile extracts the content from th input encrypted filepath (ePath),
// and creates a new decrypted file with the given key and returns the
// path of the newly created decrypted file (string) and a potential error.
// Encrypted files must end in the defined encrypted extension (eExt) to be
// accepted.
func createDFile(ePath, eExt string, key []byte) (string, error) {
	// Check to see that the given file ends in the proper file extension
	if ePath[len(ePath)-len(eExt):] != eExt {
		return "", ErrInvalidExtension
	}

	// Extract encrypted data from given file
	eData, err := os.ReadFile(ePath)
	if err != nil {
		return "", err
	}

	// Decrypt encrypted data
	dData, err := cryptography.Decrypt(eData, key)
	if err != nil {
		return "", err
	}

	dPath := genDPath(ePath, eExt)

	// Write new decrypted file
	return dPath, os.WriteFile(dPath, dData, 0644)
}

// genEPath generates the filepath of the created encrypted file so that
// it is in the same directory as the unencrypted filepath (uPath). The
// filename will end with the user-defined extension for encrypted files
// (eExt).
func genEPath(uPath, eExt string) string {
	var ePath = uPath + eExt
	return ePath
}

// genDPath generates the filepath of the created decrypted file so that
// it is in the same directory as the encrypted filepath (ePath). The new
// filwname will be stripped of the file extension for encrypted files (eExt).
func genDPath(ePath, eExt string) string {
	eName := path.Base(ePath)
	dName := "DECRYPTED-" + eName[:len(eName)-len(eExt)]
	dPath := ePath[:len(ePath)-len(eName)] + dName
	return dPath
}

// readKey acquires the keypath from the given environment variable (envVar)
// and returns the bytes of the key while chacking for errors. The key file must
// end with the file extension designated for encryption keys (keyExt)
func readKey(envVar, keyExt string) ([]byte, error) {
	// Acquire keypath from environment variable
	keyPath, envExists := os.LookupEnv(envVar)

	switch {
	case !envExists:
		return nil, ErrEnvDoesNotExist
	case keyPath == "":
		return nil, ErrEnvNotSet
	case keyPath[len(keyPath)-len(keyExt):] != keyExt: // checking for correct file extension
		return nil, ErrInvalidExtension
	case keyPath != "" && keyPath[len(keyPath)-len(keyExt):] == keyExt:
		return os.ReadFile(keyPath)
	default:
		return nil, fmt.Errorf("could not obtain environment variable")
	}
}
