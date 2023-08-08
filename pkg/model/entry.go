//nolint:ireturn
package model

import (
	"github.com/vinicius73/gamer-feed/pkg/support"
)

var _ IEntry = (*Entry)(nil) // Ensure that Entry implements IEntry.

type IEntry interface {
	Text() string
	Link() string
	Tags() []string
	Source() string
	ImageURL() string
	Hash() (string, error)

	FillFrom(IEntry) IEntry
}

type Entry struct {
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Image      string   `json:"image_url"`
	Categories []string `json:"categories"`
	SourceName string   `json:"source"`
}

// Hash of entry.
func (e Entry) Hash() (string, error) {
	return support.HashSHA256(e.URL)
}

func (e Entry) Text() string {
	return e.Title
}

func (e Entry) ImageURL() string {
	return e.Image
}

func (e Entry) Link() string {
	return e.URL
}

func (e Entry) Tags() []string {
	return e.Categories
}

func (e Entry) Source() string {
	return e.SourceName
}

func (e Entry) FillFrom(input IEntry) IEntry {
	if actual, ok := input.(Entry); ok {
		return actual
	}

	e.Title = input.Text()
	e.URL = input.Link()
	e.Image = input.ImageURL()
	e.Categories = input.Tags()
	e.SourceName = input.Source()

	return e
}
