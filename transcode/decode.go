package transcode

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"regexp"
)

func Base58Decode(encoded string) ([]byte, error) {
	decoded := big.NewInt(0)
	base := big.NewInt(58)
	alphabet := make(map[byte]int64)
	for i, char := range base58Alphabet {
		alphabet[byte(char)] = int64(i)
	}

	for i := 0; i < len(encoded); i++ {
		value, ok := alphabet[encoded[i]]
		if !ok {
			return nil, fmt.Errorf("invalid character in input")
		}

		decoded.Mul(decoded, base)
		decoded.Add(decoded, big.NewInt(value))
	}

	decodedBytes := decoded.Bytes()

	for i := 0; i < len(encoded); i++ {
		if encoded[i] == base58Alphabet[0] {
			decodedBytes = append([]byte{0x00}, decodedBytes...)
		} else {
			break
		}
	}

	return decodedBytes, nil
}

func HexDecode(encoded string) ([]byte, error) {
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %w", err)
	}
	return decoded, nil
}

func HashToHexByte(stringHash string) ([]byte, error) {
	var transactionHashByte []byte
	var err error

	re := regexp.MustCompile(`(.*)(s\d+)$`)
	matches := re.FindStringSubmatch(stringHash)

	if len(matches) > 2 {
		result, err := HexDecode(matches[1])
		if err != nil {
			log.Println("Error decoding hex string: ", err)
		}

		//transactionHashByte = []byte(result) + []byte(matches[2])

		transactionHashByte = append([]byte(result), []byte(matches[2])...)

	} else {

		transactionHashByte, err = HexDecode(stringHash)
		if err != nil {
			log.Println("Error decoding hex string: ", err)
		}

	}

	return transactionHashByte, err
}

// Base64Decode decodes the Base64 input string to a regular string.
func Base64Decode(input string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(decodedData), nil
}
