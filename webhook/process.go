package webhook

import (
	"log"
	"net/http"

	"github.com/ZeraVision/zn-wallet-manager/hmac"
)

func Process(w http.ResponseWriter, r *http.Request) {

	var txn TransactionInfo
	verified, err := hmac.VerifyRequestBody(w, r, &txn)

	if err != nil {
		log.Printf("error verifying hmac signature: %v", err)
		http.Error(w, "error verifying hmac signature", http.StatusBadRequest)
		return
	}

	if !verified {
		log.Println("hmac signature verification failed")
		http.Error(w, "hmac signature verification failed", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
