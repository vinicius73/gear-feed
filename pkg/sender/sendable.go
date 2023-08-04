package sender

import (
	"strings"

	"github.com/vinicius73/gamer-feed/pkg/scraper"
)

type Sendable interface {
	Text() string
	Link() string
	Tags() []string
	Hash() (string, error)
}

type ScrapEntry struct {
	entry scraper.Entry
}

func NewScrapEntry(entry scraper.Entry) Sendable {
	return ScrapEntry{
		entry: entry,
	}
}

func (s ScrapEntry) Text() string {
	return s.entry.Title
}

func (s ScrapEntry) Link() string {
	return s.entry.Link
}

func (s ScrapEntry) Tags() []string {
	return []string{s.entry.Type}
}

func (s ScrapEntry) Hash() (string, error) {
	return s.entry.Hash()
}

func BuildMessage(entry Sendable) string {
	var builder strings.Builder

	builder.WriteString(entry.Text())
	builder.WriteString("\n")
	builder.WriteString(entry.Link())
	builder.WriteString("\n")
	for index, tag := range entry.Tags() {
		if index > 0 {
			builder.WriteString(" ")
		}

		builder.WriteString("#" + tag)
	}

	return builder.String()
}
