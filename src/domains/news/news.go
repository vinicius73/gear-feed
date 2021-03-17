package news

import (
	"crypto/md5"
	"encoding/hex"
)

// Entry news
type Entry struct {
	Title      string
	Link       string
	Image      string
	Categories []string
	Type       string
}

// Hash of entry
func (e Entry) Hash() (string, error) {
	h := md5.New()

	_, err := h.Write([]byte(e.Link))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Key of entry
func (e Entry) Key() string {
	return e.Type + ":" + e.Title
}
