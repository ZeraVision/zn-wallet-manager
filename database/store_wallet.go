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

// InsertWallet inserts a new wallet into the database
/********** IMPORTANT **********
For best practice, it's reccomended to encrypt the address, public key and private key before storing in the database and use a kms id to reference the key in the KMS system and decrypt your data.
This allows, even in the event of a database breach, the data is still secure and can only be decrypted by your KMS system upon demand.
*/
func InsertWallet(ctx context.Context, db Database, address, publicKey, privateKey string, secretKey [32]byte, kmsId *string) error {
	// TODO -- According to user's system and requirements
	//! ** Send address (optional), public (optional) and private key for encryption (ie your KMS system) prior to storage in SQL. Note sample database structure expects string type (ie base64) **
	ePrivateKey, _ := EncryptPrivateKey(privateKey, secretKey) //! SAMPLE, local encryption (better to use external system). Should be changed according to your platform for sufficient security
	ePublicKey := string(publicKey)                            // encrypt if/as required
	eAddr := string(address)                                   // encrypt if/as required

	query := `INSERT INTO wallets (address, public_key, private_key, kms_id) VALUES ($1, $2, $3, $4)`

	err := db.Exec(ctx, query, eAddr, ePublicKey, ePrivateKey, kmsId)
	if err != nil {
		return fmt.Errorf("failed to insert wallet: %w", err)
	}

	return nil
}
