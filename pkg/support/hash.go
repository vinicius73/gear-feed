package support

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(s string) (string, error) {
	hasher := sha256.New()

	_, err := hasher.Write([]byte(s))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
