package create

import (
	"errors"

	"github.com/ZeraVision/zn-wallet-manager/transcode"
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/sha3"
)

// HashPublicKey hashes a public key using the specified algorithm.
func GetWalletAddress(publicKey []byte, hashAlg HashType, keyType KeyType) ([]byte, string, error) {
	var byteAddr []byte
	var processedPublicKey []byte

	switch hashAlg {
	case BLAKE3:
		hasher := blake3.New()
		hasher.Write(publicKey)
		byteAddr = hasher.Sum(nil)
		processedPublicKey = append([]byte("c_"), publicKey...)
	case SHA3_256:
		hasher := sha3.New256()
		hasher.Write(publicKey)
		byteAddr = hasher.Sum(nil)
		processedPublicKey = append([]byte("a_"), publicKey...)
	case SHA3_512:
		hasher := sha3.New512()
		hasher.Write(publicKey)
		byteAddr = hasher.Sum(nil)
		processedPublicKey = append([]byte("b_"), publicKey...)
	default:
		return nil, "", errors.New("unsupported hash algorithm")
	}

	if keyType == ED25519 {
		processedPublicKey = append([]byte("A_"), processedPublicKey...)
	} else if keyType == ED448 {
		processedPublicKey = append([]byte("B_"), processedPublicKey...)
	} else {
		return nil, "", errors.New("unsupported key type")
	}

	return processedPublicKey, transcode.Base58Encode(byteAddr), nil
}
