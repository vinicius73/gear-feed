package scraper

import (
	"crypto/sha256"
	"encoding/hex"
)

type Entry struct {
	Title      string
	Link       string
	Image      string
	Categories []string
	Type       string
}

// Hash of entry.
func (e Entry) Hash() (string, error) {
	hasher := sha256.New()

	_, err := hasher.Write([]byte(e.Link))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Key of entry.
func (e Entry) Key() string {
	return e.Type + ":" + e.Title
}
