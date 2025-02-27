package webhook

import (
	"context"
	"fmt"

	"github.com/ZeraVision/zn-wallet-manager/database"
)

// *** Example processing of a deposit, modify for your specific use case ***
func deposit(db database.Database, txn TransactionInfo) error {
	// Process the deposit how you see fit (updating a database, sending a notification, etc.)

	// For example, you could update a the address data to not be empty, indicating a deposit has been made.
	// if you encrypt your address in db, you can encrypt txn.Address here for the match
	///////////////////////////////////////////////////
	query := `UPDATE wallets SET empty = FALSE WHERE address = $1`
	err := db.Exec(context.Background(), query, txn.Address)

	if err != nil {
		return fmt.Errorf("failed to insert wallet: %w", err)
	}
	///////////////////////////////////////////////////

	return nil
}
