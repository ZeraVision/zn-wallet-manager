package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Function to generate an HMAC signature (client-side)
// This will take any data type (interface{}) and serialize it into a JSON string for signing
func Generate(secretKey string, data interface{}) (string, error) {
	// Serialize the data to JSON (this makes it a consistent format for signing)
	serializedData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize data: %v", err)
	}

	// Create a new HMAC instance using SHA-256
	mac := hmac.New(sha256.New, []byte(secretKey))

	// Write the serialized JSON data to the HMAC instance
	mac.Write(serializedData)

	// Compute the HMAC and return it as a base64-encoded string
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// Function to verify the HMAC signature (server-side)
// It compares the provided HMAC signature with the generated HMAC for the serialized data
func Verify(secretKey string, data interface{}, providedSignature string) (bool, error) {
	// Generate the expected HMAC signature based on the serialized data
	expectedSignature, err := Generate(secretKey, data)
	if err != nil {
		return false, err
	}

	// Compare the provided signature with the expected one (both are base64-encoded)
	return providedSignature == expectedSignature, nil
}

// VerifyRequestBody reads and verifies the request body and HMAC signature
func VerifyRequestBody(w http.ResponseWriter, r *http.Request, txn interface{}) bool {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	// Extract the HMAC signature from headers
	sigHeader := r.Header.Get("X-Signature")
	if sigHeader == "" {
		log.Printf("missing signature")
		http.Error(w, "missing signature", http.StatusUnauthorized)
		return false
	}

	// Get only the hash value after "sha256="
	sigHash := strings.TrimPrefix(sigHeader, "sha256=")

	// Deserialize JSON into txn (which should be passed as a pointer)
	err = json.Unmarshal(body, txn)
	if err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return false
	}

	// Verify HMAC signature
	verified, err := Verify(os.Getenv("SHARED_SECRET"), body, sigHash)
	if err != nil {
		log.Printf("error verifying hmac signature: %v", err)
		http.Error(w, "error verifying hmac signature", http.StatusBadRequest)
		return false
	}

	if !verified {
		log.Println("hmac signature verification failed")
		http.Error(w, "hmac signature verification failed", http.StatusUnauthorized)
		return false
	}

	return true
}
