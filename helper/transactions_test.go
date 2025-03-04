package helper

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

// Tests GetWalletBalance function on actual api with provided vars for all balance returns
func TestGetTransactions(t *testing.T) {
	if os.Getenv("INDEXER_URL") == "" {
		t.Fatalf("vars don't seem to be correctly loaded")
		return
	}

	transactions, err := GetTransactions(testAddr, nil, nil, nil, nil, true, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
		return
	}

	if transactions == nil {
		t.Fatalf("expected transactions, got nil")
		return
	}
}
