package model

import (
	"github.com/vinicius73/gear-feed/pkg/support"
)

var _ IEntry = (*Entry)(nil) // Ensure that Entry implements IEntry.

type IEntry interface {
	Text() string
	Link() string
	Tags() []string
	Source() string
	ImageURL() string
	Hash() (string, error)
	HasStory() bool
	SetHasStory(bool) IEntry

	FillFrom(IEntry) IEntry
}

type Entry struct {
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Image      string   `json:"image_url"`
	Categories []string `json:"categories"`
	SourceName string   `json:"source"`
	HaveStory  bool     `json:"has_story"`
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
	return []string{e.SourceName}
}

func (e Entry) Source() string {
	return e.SourceName
}

func (e Entry) HasStory() bool {
	return e.HaveStory
}

func (e Entry) SetHasStory(hasStory bool) IEntry {
	e.HaveStory = hasStory

	return e
}

func (e Entry) FillFrom(input IEntry) IEntry {
	if actual, ok := input.(Entry); ok {
		return actual
	}

	e = Entry{
		Title:      input.Text(),
		URL:        input.Link(),
		Image:      input.ImageURL(),
		Categories: input.Tags(),
		SourceName: input.Source(),
		HaveStory:  input.HasStory(),
	}

	return e
}
