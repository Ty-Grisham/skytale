package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	cryptography "github.com/Ty-Grisham/file-encryption"
)

const (
	tempKey = "0123456789qwerty" // hardcoding the key; will be changed later
	encExt  = ".enc"             // encyrption file extension
)

func main() {
	// Add command line options
	encrypt := flag.String("e", "", "File to be encrypted")
	decrypt := flag.String("d", "", "File to be decrypted")
	outName := flag.String("n", "", "User-specified filename for created file")

	flag.Parse()

	switch {
	// Encrypt file
	case *encrypt != "":
		cFilename, err := createEncFile(*encrypt, *outName, []byte(tempKey))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Creating encrypted file: %s\n", cFilename)

	// Decrypt file
	case *decrypt != "":
		cFilename, err := createDecFile(*decrypt, *outName, []byte(tempKey))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Creating decrypted file: %s\n", cFilename)

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
func createEncFile(uFilepath, outName string, key []byte) (string, error) {
	// Extract unencrypted data from given file
	uData, err := os.ReadFile(uFilepath)
	if err != nil {
		return "", err
	}

	// Encrypt unencrypted data
	eData, err := cryptography.Encrypt(uData, key)
	if err != nil {
		return "", err
	}

	eName := genEName(path.Base(uFilepath), outName)

	// Write new encrypted file
	return eName, os.WriteFile(eName, eData, 0644)
}

// createDeFile extracts the content from th input encrypted file,
// and creates a new decrypted file with the given key and returns the
// name of the file created (string) and a potential error. Encrypted
// files must end in the defined encrypted extension to be accepted
func createDecFile(eFilepath, outName string, key []byte) (string, error) {
	// Check to see that the given file ends in the proper file extension
	if eFilepath[len(eFilepath)-len(encExt):] != encExt {
		return "", fmt.Errorf("improper file extension; should end with %q", encExt)
	}

	// Extract encrypted data from given file
	eData, err := os.ReadFile(eFilepath)
	if err != nil {
		return "", err
	}

	// Decrypt encrypted data
	dData, err := cryptography.Decrypt(eData, key)
	if err != nil {
		return "", err
	}

	dName := genDName(path.Base(eFilepath), outName)

	// Write new decrypted file
	return dName, os.WriteFile(dName, dData, 0644)
}

// genEName creates the filename of the encrypted file
func genEName(uName, outName string) string {
	var eName string
	if outName != "" {
		// Checking to see if the .enc extension was already included
		// as part of outName
		if outName[len(outName)-len(encExt):] == encExt {
			eName = outName
		} else {
			eName = outName + encExt
		}
	} else {
		eName = uName + encExt
	}
	return eName
}

// genDName creates the filename for the decrypted file
func genDName(eName, outName string) string {
	var dName string
	if outName != "" {
		dName = outName
	} else {
		dName = eName[:len(eName)-len(encExt)]
	}

	dName = "DECRYPTED-" + dName

	return dName
}
