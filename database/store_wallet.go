package database

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

// EncryptPrivateKey encrypts the private key using a secret key.
func EncryptPrivateKey(privateKey string, secretKey [32]byte) (string, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", err
	}

	encrypted := secretbox.Seal(nonce[:], []byte(privateKey), &nonce, &secretKey)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// InsertWallet inserts a new wallet into the database with an encrypted private key.
func InsertWallet(ctx context.Context, db Database, address, publicKey, privateKey string, secretKey [32]byte) error {
	encryptedKey, err := EncryptPrivateKey(privateKey, secretKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	query := `INSERT INTO wallets (address, public_key, private_key) VALUES ($1, $2, $3)`

	err = db.Exec(ctx, query, address, publicKey, encryptedKey)
	if err != nil {
		return fmt.Errorf("failed to insert wallet: %w", err)
	}

	return nil
}
