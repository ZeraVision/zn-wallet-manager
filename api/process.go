package api

import (
	"net/http"

	"github.com/ZeraVision/zn-wallet-manager/webhook"
)

func process(w http.ResponseWriter, r *http.Request) {

	// Various external preconditions can be satisifed here as per your system requirements and design (if applicable)

	// Process deposit/withdraw hook
	webhook.Process(w, r)

}
