package model

import (
	"crypto/sha256"
	"encoding/hex"
)

type Entry struct {
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Image      string   `json:"image"`
	Categories []string `json:"categories"`
	Source     string   `json:"source"`
}

// Hash of entry.
func (e Entry) Hash() (string, error) {
	hasher := sha256.New()

	_, err := hasher.Write([]byte(e.URL))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Key of entry.
func (e Entry) Key() string {
	return e.Source + ":" + e.Title
}

func (e Entry) Text() string {
	return e.Title
}

func (e Entry) Link() string {
	return e.URL
}

func (e Entry) Tags() []string {
	return []string{e.Source}
}
