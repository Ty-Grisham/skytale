package cryptography_test

import (
	"bytes"
	"testing"

	cryptography "github.com/Ty-Grisham/skytale"
)

var (
	testKey  = []byte("0123456789qwerty") // 16 byte long constant test key
	testData = []byte("this is a test")   // date to be encrypted and decrypted
)

func TestCryptography(t *testing.T) {
	// TestConstantKey will test with the predefined testKey
	t.Run("TestConstantKey", func(t *testing.T) {
		// Encrypt data
		eData, err := cryptography.Encrypt(testData, testKey)
		if err != nil {
			t.Fatal(err)
		}

		// Decrypt data
		dData, err := cryptography.Decrypt(eData, testKey)
		if err != nil {
			t.Fatal(err)
		}

		// Assert that the decrypted data matches the original test data
		if !bytes.Equal(dData, testData) {
			t.Fatal("The decrypted data does not match the original testData")
		}
	})
}
