package transcode

import (
	"log"

	"golang.org/x/crypto/sha3"
)

// SHA3_256 calculates the SHA3-256 hash of a byte slice and returns the hash as a byte slice.
func SHA3_256(input []byte) []byte {
	hasher := sha3.New256()
	_, err := hasher.Write(input)
	if err != nil {
		log.Fatal(err)
	}
	return hasher.Sum(nil)
}
