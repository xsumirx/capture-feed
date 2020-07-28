package utility

import (
	"crypto/sha256"
	"encoding/base64"
)

func GetHash(val []byte) string {
	sum := sha256.Sum256(val)
	return base64.StdEncoding.EncodeToString(sum[:])
}