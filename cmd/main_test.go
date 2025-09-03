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
	resEPath = genEPath(uPath, eExt)
	resDPath = genDPath(resEPath, eExt)
)

// TestMain facilitates the preparing, running, and cleaning of the test
// and the resulting test artifacts
func TestMain(m *testing.M) {

	fmt.Println("Running tests...")

	// Running tests
	result := m.Run()

	fmt.Println("Cleaning up...")

	// Cleaning up
	os.Remove(resEPath)
	os.Remove(resDPath)

	os.Exit(result)
}

// TestProcess tests the overall process of encrypting/decrypting files
func TestProcess(t *testing.T) {
	// Encrypt file
	ePath, err := createEFile(uPath, eExt, tmpKey)
	if err != nil {
		t.Fatal(err)
	}

	// Decrypt file
	dPath, err := createDFile(ePath, eExt, tmpKey)
	if err != nil {
		t.Fatal(err)
	}

	// Check the bytes of the created decrypted file
	assertBytesFiles(t, dPath, uPath)
}

// TestCreatePaths tests whether the paths created are as expected
func TestCreatePaths(t *testing.T) {
	assertStrings(t, resEPath, expEPath)
	assertStrings(t, resDPath, expDPath)
}

// TestCreateEFile tests the createEFile function
func TestCreateEFile(t *testing.T) {
	// NoExpErr expects no error
	t.Run("NoExpErr", func(t *testing.T) {
		_, err := createEFile(uPath, eExt, tmpKey)
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestCreateDFile tests the createDFile function
func TestCreateDFile(t *testing.T) {
	// NoExpErr expects no errors
	t.Run("NoExpErr", func(t *testing.T) {
		_, err := createDFile(resEPath, eExt, tmpKey)
		if err != nil {
			t.Fatal(err)
		}
	})

	// InvalidExtension should return an invalid
	// extension error
	t.Run("InvalidExtension", func(t *testing.T) {
		invalidEPath := "/home/usr/sensitive_info.pdf"
		_, err := createDFile(invalidEPath, eExt, tmpKey)
		assertErrors(t, err, ErrInvalidExtension)
	})
}

func TestReadKey(t *testing.T) {
	// ExpectedInput should read key and produce no errors
	t.Run("ExpectedInput", func(t *testing.T) {
		t.Setenv(envVar, testKeyPath)
		_, err := readKey(envVar, keyExt)
		if err != nil {
			t.Errorf("Error in ExpectedInput: %q", err)
		}
	})

	// EnvNotSet tests the output of readKey with no environment variable set
	// e.g. envVar=""
	t.Run("EnvNotSet", func(t *testing.T) {
		t.Setenv(envVar, "")
		_, err := readKey(envVar, keyExt)
		assertErrors(t, err, ErrEnvNotSet)
	})

	// EnvDoesNotExist tests the output of readKey when the environment variable
	// does not exist
	t.Run("EnvDoesNotExist", func(t *testing.T) {
		os.Unsetenv(envVar)
		_, err := readKey(envVar, keyExt)
		assertErrors(t, err, ErrEnvDoesNotExist)
	})

	// InvalidKeyExtension tests the output of readKey when the key file in the
	// environment variable has an invalid file extension
	t.Run("InvalidKeyExtension", func(t *testing.T) {
		t.Setenv(envVar, uPath)
		_, err := readKey(envVar, keyExt)
		assertErrors(t, err, ErrInvalidExtension)
	})
}

// assertErrors is a helper function that asserts that an error should be raised and
// and the error should be the one specified in the input
func assertErrors(t *testing.T, res, exp error) {
	t.Helper()
	if res == nil {
		t.Fatalf("Expected %s, but no errors were raised", exp)
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
func assertBytesFiles(t *testing.T, resFile, expFile string) {
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
