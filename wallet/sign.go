package wallet

import (
	"errors"
	"fmt"

	"github.com/GoKillers/libsodium-go/cryptosign"
	"github.com/ZeraVision/zn-wallet-manager/transcode"
	"github.com/cloudflare/circl/sign/ed448"
)

// SignTransaction signs a transaction payload using the Ed25519 private key.
func Sign(privateKeyBase58 string, payload []byte, keyType KeyType) ([]byte, error) {
	if len(payload) == 0 {
		return nil, errors.New("payload cannot be empty")
	}

	// Decode the Base58-encoded private key
	privateKey, err := transcode.Base58Decode(privateKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	// Ensure the private key length matches Ed25519 requirements
	if len(privateKey) != cryptosign.CryptoSignSecretKeyBytes() {
		return nil, errors.New("invalid private key length")
	}

	var signature []byte

	if keyType == ED25519 {
		var retCode int
		signature, retCode = cryptosign.CryptoSignDetached(payload, privateKey)
		if retCode != 0 {
			return nil, errors.New("libsodium: failed to sign transaction")
		}
	} else if keyType == ED448 {
		privKey := ed448.NewKeyFromSeed(privateKey)
		signature := ed448.Sign(privKey, payload, "")
		return signature, nil
	} else {
		return nil, errors.New("unsupported key type")
	}

	return signature, nil
}
