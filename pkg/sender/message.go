package sender

import (
	"strings"

	"github.com/vinicius73/gamer-feed/pkg/model"
)

func BuildMessage(entry model.IEntry) string {
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
