package news

import (
	"crypto/md5"
	"encoding/hex"
)

// Entry news
type Entry struct {
	Title string
	Link  string
	Image string
	Type  string
}

// Hash of entry
func (e Entry) Hash() string {
	h := md5.New()

	h.Write([]byte(e.Link))

	return hex.EncodeToString(h.Sum(nil))
}

// Key of entry
func (e Entry) Key() string {
	return e.Type + ":" + e.Title
}
