package webhook

import "math/big"

type depositWithdraw string

const (
	Deposit  = "deposit"
	Withdraw = "withdrawl"
)

type TransactionInfo struct {
	Address        string          `json:"address"`
	Type           depositWithdraw `json:"type"`
	Hash           string          `json:"hash"`
	Symbol         string          `json:"symbol"`
	AmountParts    *big.Int        `json:"amountParts"`
	PartsPerSymbol *big.Int        `json:"partsPerSymbol"`
	Amount         *big.Float      `json:"amount"`
}
