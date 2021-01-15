package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var Base64URLNoPadding = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString

// Generate a random Base64 string that is url query safe.
func GenerateBase64URLnopadding(length uint32) string {
	randBytes := make([]byte, length)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	return Base64URLNoPadding(randBytes)
}

// Generate SHA256 hash.
func GenerateSHA256(input string) [32]byte {
	return sha256.Sum256([]byte(input))
}

// Encrypt a byte slice with a key.
//
// Source:
// https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
func Encrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) > 32 {
		return nil, errors.New("Key must be less than equal to 32 bytes!")
	}

	// pad key to 32 bytes
	key = []byte(fmt.Sprintf("%032v", string(key)))

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt a byte slice with a key.
//
// Source:
// https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
func Decrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) > 32 {
		return nil, errors.New("Key must be less than equal to 32 bytes!")
	}

	// pad key to 32 bytes
	key = []byte(fmt.Sprintf("%032v", string(key)))

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	noncesize := gcm.NonceSize()
	if len(data) < noncesize {
		return nil, errors.New("File is not valid")
	}

	nonce, data := data[:noncesize], data[noncesize:]

	plaindata, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return plaindata, nil
}
