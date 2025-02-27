package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/nacl/secretbox"
)

// EncryptedKeyMatcher implements sqlmock.Argument to verify the encrypted private key.
type EncryptedKeyMatcher struct {
	privateKey string
	secretKey  [32]byte
}

// Match checks if the provided value (as driver.Value) is a valid encrypted key that decrypts to the expected private key.
func (m EncryptedKeyMatcher) Match(v driver.Value) bool {
	var str string
	switch value := v.(type) {
	case string:
		str = value
	case []byte:
		str = string(value)
	default:
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return false
	}

	if len(decoded) < 24 {
		return false
	}

	var nonce [24]byte
	copy(nonce[:], decoded[:24])
	decrypted, ok := secretbox.Open(nil, decoded[24:], &nonce, &m.secretKey)
	return ok && string(decrypted) == m.privateKey
}

// MockDB wraps a sql.DB instance to implement the Database interface.
type MockDB struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

// Exec calls the underlying ExecContext on the sql.DB.
func (m *MockDB) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

func TestInsertWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	mockDB := &MockDB{db: db, mock: mock}

	address := "test_address"
	publicKey := "test_public_key"
	privateKey := "test_private_key"
	kmsID := func(s string) *string { return &s }("test_kms_id")

	var secretKey [32]byte
	if _, err := rand.Read(secretKey[:]); err != nil {
		t.Fatalf("failed to generate secret key: %v", err)
	}

	// Create an instance of EncryptedKeyMatcher
	privateKeyMatcher := EncryptedKeyMatcher{privateKey: privateKey, secretKey: secretKey}

	// Set expectation using the custom matcher for the encrypted private key.
	mock.ExpectExec(`INSERT INTO wallets \(address, public_key, private_key, kms_id\) VALUES \(\$1, \$2, \$3, \$4\)`).
		WithArgs(address, publicKey, privateKeyMatcher, kmsID). // Pass matcher instance, not inline struct
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = InsertWallet(ctx, mockDB, address, publicKey, privateKey, secretKey, kmsID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestSelectAndDecryptWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	address := "test_address"
	publicKey := "test_public_key"
	privateKey := "test_private_key"
	var secretKey [32]byte
	if _, err := rand.Read(secretKey[:]); err != nil {
		t.Fatalf("failed to generate secret key: %v", err)
	}

	// Encrypt the private key using the production function.
	encryptedKey, err := EncryptPrivateKey(privateKey, secretKey)
	if err != nil {
		t.Fatalf("failed to encrypt private key: %v", err)
	}

	// Prepare the row that will be returned by the SELECT query.
	rows := sqlmock.NewRows([]string{"address", "public_key", "private_key"}).
		AddRow(address, publicKey, encryptedKey)

	// Set expectation for the SELECT query.
	mock.ExpectQuery(`SELECT address, public_key, private_key FROM wallets WHERE address = \$1`).
		WithArgs(address).
		WillReturnRows(rows)

	// Execute the query.
	var addr, pub, enc string
	err = db.QueryRow("SELECT address, public_key, private_key FROM wallets WHERE address = $1", address).
		Scan(&addr, &pub, &enc)
	if err != nil {
		t.Fatalf("failed to query row: %v", err)
	}

	// Decrypt the encrypted private key.
	decoded, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		t.Fatalf("failed to decode encrypted key: %v", err)
	}
	if len(decoded) < 24 {
		t.Fatalf("decoded data is too short to contain a nonce")
	}
	var nonce [24]byte
	copy(nonce[:], decoded[:24])
	decrypted, ok := secretbox.Open(nil, decoded[24:], &nonce, &secretKey)
	if !ok {
		t.Fatalf("failed to decrypt private key")
	}
	if string(decrypted) != privateKey {
		t.Fatalf("expected private key %s, got %s", privateKey, string(decrypted))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}
