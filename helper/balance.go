package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
)

type Token struct {
	Symbol        string   `json:"symbol"`
	Name          string   `json:"name"`
	Icon          *string  `json:"icon,omitempty"`
	Amount        *big.Int `json:"amount"`
	Parts         *big.Int `json:"parts"`
	Rate          float64  `json:"rate"`
	Value         float64  `json:"value"`
	ChangePercent float64  `json:"changePercent"`
	Type          string   `json:"type"`
}

type WalletResponse struct {
	TotalWalletValue float64 `json:"totalWalletValue"`
	ZraRate          float64 `json:"zraRate"`
	Tokens           struct {
		Found  int     `json:"found"`
		Tokens []Token `json:"tokens"`
	} `json:"tokens"`
}

// Get all by passing symbol as nil
func GetWalletBalance(address string, symbol *string) (*WalletResponse, error) {
	baseUrl := fmt.Sprintf("%s/store?requestType=getWalletBalance&address=%s&filter=tokens", os.Getenv("INDEXER_URL"), address)

	if symbol != nil {
		baseUrl += fmt.Sprintf("&symbol=%s", *symbol)
	}

	limit := 100 // default max tokens per call
	if symbol != nil {
		limit = 1
	}
	baseUrl += fmt.Sprintf("&limit=%d", limit)

	var combinedWalletResponse *WalletResponse
	offset := 0

	for {
		var walletResponse WalletResponse
		url := baseUrl + fmt.Sprintf("&offset=%d", offset)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Target", "explorer")
		req.Header.Set("Authorization", "Api-Key "+os.Getenv("INDEXER_API_KEY"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
		}
		if combinedWalletResponse == nil {
			combinedWalletResponse = new(WalletResponse)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(body, &walletResponse); err != nil {
			return nil, err
		}

		if offset == 0 {
			*combinedWalletResponse = walletResponse
		} else {
			combinedWalletResponse.Tokens.Tokens = append(combinedWalletResponse.Tokens.Tokens, walletResponse.Tokens.Tokens...)
		}

		count := len(walletResponse.Tokens.Tokens)
		offset += count
		if combinedWalletResponse.Tokens.Found == offset { // Stop if no more data
			break
		} else if count == 0 {
			return combinedWalletResponse, fmt.Errorf("expected 0 tokens found in search")
		}

	}

	return combinedWalletResponse, nil
}
