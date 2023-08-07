package model

import (
	"github.com/vinicius73/gamer-feed/pkg/support"
)

var _ IEntry = (*Entry)(nil) // Ensure that Entry implements IEntry.

type IEntry interface {
	Text() string
	Link() string
	Tags() []string
	Hash() (string, error)
}

type Entry struct {
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Image      string   `json:"image"`
	Categories []string `json:"categories"`
	Source     string   `json:"source"`
}

// Hash of entry.
func (e Entry) Hash() (string, error) {
	return support.HashSHA256(e.URL)
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
