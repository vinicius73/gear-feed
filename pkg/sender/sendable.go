package sender

import (
	"strings"

	"github.com/vinicius73/gamer-feed/pkg/model"
)

var _ Sendable = model.Entry{} // Ensure that Entry implements Sendable.

type Sendable interface {
	Text() string
	Link() string
	Tags() []string
	Hash() (string, error)
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
