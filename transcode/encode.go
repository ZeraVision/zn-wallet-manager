package transcode

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func Base58EncodePublicKey(publicKey []byte) string {
	publicKeyParts := bytes.SplitN(publicKey, []byte{'_'}, 2)
	if len(publicKeyParts) != 2 {
		return "invalid_key"
	}
	return string(publicKeyParts[0]) + Base58Encode(publicKeyParts[1])
}

func Base58Encode(input interface{}) string {
	var data []byte
	switch v := input.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return ""
	}

	x := big.NewInt(0).SetBytes(data)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	encoded := ""

	for x.Cmp(zero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, base, mod)
		encoded = string(base58Alphabet[mod.Int64()]) + encoded
	}

	for _, b := range data {
		if b != 0x00 {
			break
		}
		encoded = string(base58Alphabet[0]) + encoded
	}

	return encoded
}

func HexEncode(input []byte) string {
	return hex.EncodeToString(input)
}

// Base64Encode encodes the input byte slice to Base64.
func Base64Encode(input []byte) string {
	encodedData := base64.StdEncoding.EncodeToString(input)
	return encodedData
}
