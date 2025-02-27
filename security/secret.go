package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

// GenerateSecretKey creates a new 256-bit secret key and returns it in Base64 format.
func GenerateSecretKey() (string, error) {
	var key [32]byte
	if _, err := rand.Read(key[:]); err != nil {
		return "", fmt.Errorf("failed to generate secret key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(key[:]), nil
}

// DecodeSecretKey converts a Base64-encoded key to a 32-byte array.
func DecodeSecretKey(encodedKey string) ([32]byte, error) {
	var key [32]byte
	decoded, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return key, fmt.Errorf("failed to decode secret key: %w", err)
	}

	if len(decoded) != 32 {
		return key, errors.New("invalid secret key length")
	}

	copy(key[:], decoded)
	return key, nil
}

// DecryptPrivateKey decrypts an encrypted private key using the given secret key.
func DecryptPrivateKey(encryptedKey string, secretKey [32]byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted key: %w", err)
	}

	if len(decoded) < 24 {
		return "", errors.New("encrypted key is too short")
	}

	var nonce [24]byte
	copy(nonce[:], decoded[:24])
	decrypted, ok := secretbox.Open(nil, decoded[24:], &nonce, &secretKey)
	if !ok {
		return "", errors.New("failed to decrypt private key")
	}
	return string(decrypted), nil
}
