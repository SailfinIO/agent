package agent

import (
	"crypto/rand"
	"encoding/hex"
)

// generateAPIKey generates a new 32-byte secure random API key encoded as a hex string.
func generateAPIKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
