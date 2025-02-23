package sender

import (
	"strconv"
	"strings"
	"time"

	"github.com/vinicius73/gear-feed/pkg"
	"github.com/vinicius73/gear-feed/pkg/model"
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

func BuildCleanupMessage(count int64) string {
	var builder strings.Builder

	builder.WriteString(BuildMsgHeader())

	builder.WriteRune('\n')

	builder.WriteString("🧹 <b>Cleanup: </b><code>")
	builder.WriteString(strconv.FormatInt(count, 10))
	builder.WriteString(" records</code>")

	builder.WriteString(BuildMsgFooter())

	return builder.String()
}

func BuildMsgHeader() string {
	var builder strings.Builder

	builder.WriteString("ℹ️ <b>")
	builder.WriteString(pkg.AppName)
	builder.WriteString(" - <code>")
	builder.WriteString(pkg.Host())
	builder.WriteString("</code></b>\n🤖 <i>")
	builder.WriteString(pkg.Version())
	builder.WriteRune(' ')
	builder.WriteString(pkg.Commit())
	builder.WriteString("</i>")

	return builder.String()
}

func BuildMsgFooter() string {
	var builder strings.Builder

	builder.WriteString("\n\n 🕔 <i>")
	builder.WriteString(time.Now().Format(time.RFC3339))
	builder.WriteString("</i>")

	return builder.String()
}
