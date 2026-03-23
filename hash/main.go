package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func hashRowKey(rowKey string) string {
	// Convert row key to bytes
	rowKeyBytes := []byte(rowKey)

	// Hash the row key using SHA-256
	hash := sha256.Sum256(rowKeyBytes)

	// Convert the hash to a hexadecimal string
	hashedKey := hex.EncodeToString(hash[:])

	// Return the hashed key
	return hashedKey
}

func main() {
	// Example usage
	rowKey := "1.187.0.0/25"
	hashedKey := hashRowKey(rowKey)
	fmt.Println("Hashed Row Key:", hashedKey)
}
