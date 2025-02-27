//! Integration test based on actual api calls to test env setup

package helper

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var testAddr = "GLoZ2hgtqUBnjoPkDQ4pYbVbJ7CLmvfUSrrmGpGq7zJK"

func init() {
	godotenv.Load("../.env")

	//! if needed for your environment, set your sample keys here... (bad)
	// os.Setenv("SAMPLE_INDEXER_URL", os.Getenv("INDEXER_URL"))
	// os.Setenv("INDEXER_API_KEY", os.Getenv("INDEXER_API_KEY"))
}

// Tests GetWalletBalance function on actual api with provided vars for all balance returns
func TestGetAllWalletBalance(t *testing.T) {
	if os.Getenv("INDEXER_URL") == "" {
		t.Fatalf("vars don't seem to be correctly loaded")
		return
	}

	walletResponse, err := GetWalletBalance(testAddr, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
		return
	}

	if walletResponse == nil {
		t.Fatalf("expected walletResponse, got nil")
		return
	}
}

// Tests GetWalletBalance function on actual api with provided vars for single balance return
func TestGetAllWalletBalanceSingle(t *testing.T) {
	if os.Getenv("INDEXER_URL") == "" {
		t.Fatalf("vars don't seem to be correctly loaded")
		return
	}

	symbol := "$ZRA+0000"

	walletResponse, err := GetWalletBalance(testAddr, &symbol)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
		return
	}

	if walletResponse == nil {
		t.Fatalf("expected walletResponse, got nil")
		return
	}

	if walletResponse.Tokens.Found != 1 {
		t.Fatalf("expected 1 token, got %v", walletResponse.Tokens.Found)
		return
	}

	if walletResponse.Tokens.Tokens[0].Symbol != symbol {
		t.Fatalf("expected %v, got %v", symbol, walletResponse.Tokens.Tokens[0].Symbol)
		return
	}
}
