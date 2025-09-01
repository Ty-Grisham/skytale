package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// encrypt encrypts the input data with the input key using AES-GCM.
// Returns a slice of bytes containing encrypted data or a potential error
func Encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil // first argument must be nonce, contrary to crypto/cipher docs
}

// decrypt decrypts the input encrypted data with the input key using AES-GCM.
// Returns a slice of bytes containing the decrypted data or a potential error
func Decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("data to be decrypted is too short")
	}

	nonce, data := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, data, nil)
}
