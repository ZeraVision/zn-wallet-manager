package webhook

import (
	"net/http"

	"github.com/ZeraVision/zn-wallet-manager/database"
	"github.com/ZeraVision/zn-wallet-manager/hmac"
)

func Process(w http.ResponseWriter, r *http.Request) {

	db := database.Get()

	var txn TransactionInfo
	// If signature doesn't match, reject
	if !hmac.VerifyRequestBody(w, r, &txn) {
		return
	}

	// On successful parse and verification, respond with 200 OK
	w.WriteHeader(http.StatusOK)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// TODO update these functions to include your required logic
	if txn.Type == Deposit {
		deposit(db, txn)
	} else if txn.Type == Withdraw {
		withdrawl(db, txn)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
