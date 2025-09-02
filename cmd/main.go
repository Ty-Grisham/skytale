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
	envVar = "AES_KEY" //  Name of environment variable for encryption key
	eExt   = ".enc"    // encyrption file extension
)

func main() {
	// Add command line options
	encrypt := flag.String("e", "", "File to be encrypted")
	decrypt := flag.String("d", "", "File to be decrypted")

	flag.Parse()

	// Acquire key from environment variable
	key, err := readKey(envVar)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	switch {
	// Encrypt file
	case *encrypt != "":
		cFilename, err := createEncFile(*encrypt, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Created encrypted file: %s\n", cFilename)

	// Decrypt file
	case *decrypt != "":
		cFilename, err := createDecFile(*decrypt, key)
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

// createEncFile extracts the content from the input unencrypted file,
// and creates a new encrypted file with the given key and returns the
// name of the created encrypted file (string) a potential error
func createEncFile(uPath string, key []byte) (string, error) {
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

	eName := genEPath(uPath)

	// Write new encrypted file
	return eName, os.WriteFile(eName, eData, 0644)
}

// createDeFile extracts the content from th input encrypted file,
// and creates a new decrypted file with the given key and returns the
// name of the file created (string) and a potential error. Encrypted
// files must end in the defined encrypted extension to be accepted
func createDecFile(ePath string, key []byte) (string, error) {
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

	dPath := genDPath(ePath)

	// Write new decrypted file
	return dPath, os.WriteFile(dPath, dData, 0644)
}

// genEPath generates the filepath of the created file encrypted file so that
// it is in the same directory as the uPath (unencrypted filepath)
func genEPath(uPath string) string {
	var ePath = uPath + ".enc"
	// if outName != "" {
	// 	// Checking to see if the .enc extension was already included
	// 	// as part of outName
	// 	if outName[len(outName)-len(encExt):] == encExt {
	// 		eName = outName
	// 	} else {
	// 		eName = outName + encExt
	// 	}
	// } else {
	// 	ePath = uName + encExt
	// }

	return ePath
}

// genDName creates the filename for the decrypted file
func genDPath(ePath string) string {
	eName := path.Base(ePath)
	dName := "DECRYPTED-" + eName[:len(eName)-len(eExt)]
	dPath := ePath[:len(ePath)-len(eName)] + dName
	// if outName != "" {
	// 	dName = outName
	// } else {
	// 	dName = eName[:len(eName)-len(encExt)]
	// }

	return dPath
}

// getKey
func readKey(envVar string) ([]byte, error) {
	// Acquire keypath from environment variable
	keyPath, envExists := os.LookupEnv(envVar)

	switch {
	case !envExists:
		return nil, ErrEnvDoesNotExist
	case keyPath == "":
		return nil, ErrEnvNotSet
	case keyPath != "":
		return os.ReadFile(keyPath)
	default:
		return nil, fmt.Errorf("could not obtain environment variable")
	}
}
