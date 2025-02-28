package wallet

import (
	"math/big"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestEd25519(t *testing.T) {
	testSignature(t, "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR", "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7", "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs", ED25519)
}

func TestEd448(t *testing.T) {
	testSignature(t, "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS", "B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV", "HYkGjJY8hjEAxLe1UFzEni5mANwbvTquvTV6mgMT6Qp2Ee1CFYC8tVNfdqyJ9ZwnwsYRUwfMg15suW", ED448)
}

func testSignature(t *testing.T, address, testPublic, testPrivate string, keyType KeyType) {

	inputs := []Inputs{}

	inputs = append(inputs, Inputs{
		B58Address: address,
		KeyType:    keyType,
		PublicKey:  testPublic,
		PrivateKey: testPrivate,
		Amount:     big.NewFloat(1.01),
		FeePercent: 100,
	})

	outputs := map[string]*big.Float{}

	outputs["outputAddr1"] = big.NewFloat(1.01)

	symbol := "$ZRA+0000"
	baseFeeID := "$ZRA+0000"
	baseFeeAmountParts := "1000000000" // 1 zra

	txn, err := CreateCoinTxn(inputs, outputs, symbol, baseFeeID, baseFeeAmountParts, nil, nil)

	if err != nil {
		t.Errorf("Error creating transaction: %s", err)
	}

	// Grab signature
	signature := txn.Auth.Signature[0]

	// Remove signature & hash before verification
	txn.Auth.Signature = nil
	txn.Base.Hash = nil

	txnBytes, err := proto.Marshal(txn)
	if err != nil {
		t.Errorf("Error marshalling transaction: %s", err)
	}

	// Verify signature of object
	// TODO here

	ok, err := Verify(testPublic, txnBytes, signature)

	if err != nil {
		t.Errorf("Error verifying signature: %s", err)
	}

	if !ok {
		t.Errorf("Signature verification failed")
	}
}
