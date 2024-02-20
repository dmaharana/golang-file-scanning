package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const DefaultKey = "2320876f6ce0420ea20258b5fc54669b9af3babf"

// createHash  generates a hash from the given string. It is used to generate unique keys for AES encryption.
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// encryptWithCrypto returns  the encrypted version of data using aes cbc encryption with the provided key and iv.
func encryptWithCrypto(data []byte, passphrase string) ([]byte, error) {
	if passphrase == "" {
		passphrase = DefaultKey
	}

	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// decryptWithCrypto  returns the original data from an encrypted slice and any error that occurred during decryption
func decryptWithCrypto(data []byte, passphrase string) ([]byte, error) {
	if passphrase == "" {
		passphrase = DefaultKey
	}

	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Encrypt  takes the given data and an optional passphrase to encrypt it with AES-256
func Encrypt(text string, key string) (string, bool) {
	success := false
	if IsEncrypted(text) {
		return text, success
	}

	encoded, err := encryptWithCrypto([]byte(text), key)
	if err != nil {
		return text, success
	}

	hEncoded := hex.EncodeToString(encoded)

	success = true
	return hEncoded, success
}

// Decrypt  decrypts the given encrypted text using the provided passphrase and returns it as a string.
func Decrypt(text string, key string) (string, bool) {
	success := false

	if !IsEncrypted(text) {
		return text, success
	}

	decoded, err := hex.DecodeString(text)
	if err != nil {
		return text, success
	}

	decoded, err = decryptWithCrypto(decoded, key)
	if err != nil {
		return text, success
	}

	success = true
	return string(decoded), success

}

// IsEncrypted  returns whether or not the given cipherText is encrypted.
func IsEncrypted(txt string) bool {
	encrypted := true

	_, err := hex.DecodeString(txt)
	if err != nil {
		encrypted = false
	}

	return encrypted
}
