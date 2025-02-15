package security

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/nacl/secretbox"
)

func TestGenerateSecretKey(t *testing.T) {
	key, err := GenerateSecretKey()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(key) == 0 {
		t.Fatalf("expected non-empty key, got empty key")
	}

	decodedKey, err := DecodeSecretKey(key)
	if err != nil {
		t.Fatalf("expected no error decoding key, got %v", err)
	}

	if len(decodedKey) != 32 {
		t.Fatalf("expected key length of 32, got %d", len(decodedKey))
	}
}

func TestDecodeSecretKey(t *testing.T) {
	validKey := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	invalidKey := "invalid_base64"

	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"ValidKey", validKey, false},
		{"InvalidBase64", invalidKey, true},
		{"InvalidLength", "short", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeSecretKey(tt.input)
			if tt.expectErr && err == nil {
				t.Fatalf("expected error, got none")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
func TestDecryptPrivateKey(t *testing.T) {
	secretKey := [32]byte{}
	if _, err := rand.Read(secretKey[:]); err != nil {
		t.Fatalf("failed to generate secret key: %v", err)
	}

	privateKey := "test_private_key"
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		t.Fatalf("failed to generate nonce: %v", err)
	}

	encrypted := secretbox.Seal(nonce[:], []byte(privateKey), &nonce, &secretKey)
	encryptedKey := base64.StdEncoding.EncodeToString(encrypted)

	decryptedKey, err := DecryptPrivateKey(encryptedKey, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if decryptedKey != privateKey {
		t.Fatalf("expected decrypted key to be %s, got %s", privateKey, decryptedKey)
	}
}

func TestDecryptPrivateKey_InvalidBase64(t *testing.T) {
	secretKey := [32]byte{}
	invalidBase64 := "invalid_base64"

	_, err := DecryptPrivateKey(invalidBase64, secretKey)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestDecryptPrivateKey_TooShort(t *testing.T) {
	secretKey := [32]byte{}
	tooShort := base64.StdEncoding.EncodeToString([]byte("short"))

	_, err := DecryptPrivateKey(tooShort, secretKey)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestDecryptPrivateKey_FailedDecryption(t *testing.T) {
	secretKey := [32]byte{}
	if _, err := rand.Read(secretKey[:]); err != nil {
		t.Fatalf("failed to generate secret key: %v", err)
	}

	privateKey := "test_private_key"
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		t.Fatalf("failed to generate nonce: %v", err)
	}

	encrypted := secretbox.Seal(nonce[:], []byte(privateKey), &nonce, &secretKey)
	encryptedKey := base64.StdEncoding.EncodeToString(encrypted)

	// Modify the encrypted key to simulate decryption failure
	encryptedKey = encryptedKey[:len(encryptedKey)-1] + "A"

	_, err := DecryptPrivateKey(encryptedKey, secretKey)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}
