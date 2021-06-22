package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(data string) string {
	dt := []byte(data)
	hash := sha256.Sum256(dt)
	return hex.EncodeToString(hash[:])
}