package main

import (
	"bytes"
	"os"
	"path"
	"testing"
)

const (
	uFilepath = "../testdata/test1.md"
	expEName  = "test1.md.enc"
	expDName  = "DECRYPTED-test1.md"
)

var (
	resEName = genEName(path.Base(uFilepath), "")
	resDName = genDName(path.Base(resEName), "")
)

// TestMain facilitates the preparing, running, and cleaning of the test
// and the resulting test artifacts
func TestMain(m *testing.M) {
	// Running tests
	result := m.Run()

	// Cleaning up
	os.Remove(resEName)
	os.Remove(resDName)

	os.Exit(result)
}

// TestExpexted will test that all phases of the program will work
// given the proper expected inputs and no specified outName
func TestFunctionality(t *testing.T) {
	// TestBasicFunctionality tests the basic encrypting/decrypting of the
	// test file
	t.Run("TestBasic", func(t *testing.T) {
		// Encrypting file
		eFilename, err := createEncFile(uFilepath, "", []byte(tempKey))
		if err != nil {
			t.Fatal(err)
		}

		// Decrypting file
		dFilename, err := createDecFile(eFilename, "", []byte(tempKey))
		if err != nil {
			t.Fatal(err)
		}

		// Checking that the file was properly encrypted and decrypted
		assertBytesFiles(t, uFilepath, dFilename)
	})

	// TestLogistics tests whether the created names are as expected
	t.Run("TestLogistics", func(t *testing.T) {
		assertStrings(t, expEName, resEName)
		assertStrings(t, expDName, resDName)
	})
}

// assertStrings is a helper function that compares the input result string
// to the input expected dtring
func assertStrings(t *testing.T, exp, res string) {
	t.Helper()

	// Compare strings
	if res != exp {
		t.Fatalf("Expected %q, got %q\n", exp, res)
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
		t.Fatal("resulting file does not match the expected file")
	}

}
