package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

const (
	testKeyPath = "../testdata/testKey.aes"
	uPath       = "../testdata/test1.md"
	expEPath    = "../testdata/test1.md.enc"
	expDPath    = "../testdata/DECRYPTED-test1.md"
)

var (
	tmpKey   = []byte("0123456789qwerty")
	resEPath = genEPath(uPath)
	resDPath = genDPath(resEPath)
)

// TestMain facilitates the preparing, running, and cleaning of the test
// and the resulting test artifacts
func TestMain(m *testing.M) {

	fmt.Println("Setting up environment...")

	// Setting up environment
	ogEnv := os.Getenv(envVar)

	fmt.Println("Running tests...")

	// Running tests
	result := m.Run()

	fmt.Println("Cleaning up...")

	// Cleaning up
	os.Remove(resEPath)
	os.Remove(resDPath)
	os.Setenv(envVar, ogEnv)

	os.Exit(result)
}

// TestExpexted will test that all phases of the program will work
// given the proper expected inputs and no specified outName
func TestFunctionality(t *testing.T) {
	// TestBasicFunctionality tests the basic encrypting/decrypting of the
	// test file
	t.Run("TestBasic", func(t *testing.T) {
		// Encrypting file
		ePath, err := createEncFile(uPath, []byte(tmpKey))
		if err != nil {
			t.Fatal(err)
		}

		// Decrypting file
		dPath, err := createDecFile(ePath, []byte(tmpKey))
		if err != nil {
			t.Fatal(err)
		}

		// Checking that the file was properly encrypted and decrypted
		assertBytesFiles(t, uPath, dPath)
	})

	// TestLogistics tests whether the created names are as expected
	t.Run("FileNaming", func(t *testing.T) {
		assertStrings(t, expEPath, resEPath)
		assertStrings(t, expDPath, resDPath)
	})
}

// TestErrors tests the program with input that should raise errors
func TestErrors(t *testing.T) {
	// InvalidExtension tests the createDecFile function when it is provided with
	// an input filepath that has an invalid file extension
	t.Run("InvalidExtension", func(t *testing.T) {
		// Attempting to decrypt file
		dPath, err := createDecFile(uPath, []byte(tmpKey)) // The extension of uPath is .md
		assertErrors(t, err, ErrInvalidExtension)
		os.Remove(dPath) // File is deleted if created
	})
}

func TestReadKey(t *testing.T) {
	// ExpectedInput should read key and produce no errors
	t.Run("ExpectedInput", func(t *testing.T) {
		os.Setenv(envVar, testKeyPath)
		_, err := readKey(envVar)
		if err != nil {
			t.Errorf("Error in ExpectedInput: %q", err)
		}
	})

	// EnvNotSet tests the output of readKey with no environment variable set
	// e.g. envVar=""
	t.Run("EnvNotSet", func(t *testing.T) {
		os.Setenv(envVar, "")
		_, err := readKey(envVar)
		assertErrors(t, err, ErrEnvNotSet)
	})

	// EnvDoesNotExist tests the output of readKey when the environment variable
	// does not exist
	t.Run("EnvDoesNotExist", func(t *testing.T) {
		os.Unsetenv(envVar)
		_, err := readKey(envVar)
		assertErrors(t, err, ErrEnvDoesNotExist)
	})
}

// assertErrors is a helper function that asserts that an error should be raised and
// and the error should be the one specified in the input
func assertErrors(t *testing.T, res, exp error) {
	t.Helper()
	if res == nil {
		t.Fatal("Expected error, but none were raised")
	}

	if res != exp {
		t.Errorf("Expected %q, got %q", exp, res)
	}
}

// assertStrings is a helper function that compares the input result string
// to the input expected dtring
func assertStrings(t *testing.T, res, exp string) {
	t.Helper()

	// Compare strings
	if res != exp {
		t.Errorf("Expected %q, got %q", exp, res)
	}
}

// assertBytes is a helper function that checks that bytes in the result file
// are the same as the bytes in the expected file
func assertBytesFiles(t *testing.T, expFile, resFile string) {
	t.Helper()

	// Read data from files
	res, err := os.ReadFile(resFile)
	if err != nil {
		t.Fatal(err)
	}

	exp, err := os.ReadFile(expFile)
	if err != nil {
		t.Fatal(err)
	}

	// Compare bytes
	if !bytes.Equal(res, exp) {
		t.Error("resulting file does not match the expected file")
	}

}
