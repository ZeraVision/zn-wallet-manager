package wallet

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Gets the nonce from ZV indexer, ignores ssl cert in case its self signed or a non-official service to get nonce
func GetNonce(address string) uint64 {
	// Create a custom transport with TLSClientConfig to skip certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}
	url := fmt.Sprintf(os.Getenv("INDEXER_URL")+"/store?requestType=getNextNonce&address=%s", address)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	req.Header.Add("Target", "explorer")
	req.Header.Add("authorization", "Api-Key "+os.Getenv("INDEXER_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	nonce, err := strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return nonce
}
