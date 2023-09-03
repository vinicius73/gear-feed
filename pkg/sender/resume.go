package sender

import (
	"strconv"
	"strings"
)

type Resume struct {
	Loaded   int
	Filtered int
	Sources  []ResumeSource
}

type ResumeSource struct {
	Source   string
	Loaded   int
	Filtered int
}

func (r Resume) HTML() string {
	var builder strings.Builder

	builder.WriteString(BuildMsgHeader())

	builder.WriteRune('\n')
	builder.WriteString("🗞 <b>Resume: </b>")
	builder.WriteString("<code>")
	builder.WriteString(strconv.Itoa(r.Loaded))
	builder.WriteRune('/')
	builder.WriteString(strconv.Itoa(r.Filtered))
	builder.WriteString("</code>")
	builder.WriteRune('\n')

	for _, source := range r.Sources {
		builder.WriteRune('\n')
		builder.WriteString(source.HTML())
	}

	builder.WriteString(BuildMsgFooter())

	return builder.String()
}

func (r ResumeSource) HTML() string {
	var builder strings.Builder

	builder.WriteString("- <b>")
	builder.WriteString(r.Source)
	builder.WriteString("</b>: <code>")
	builder.WriteString(strconv.Itoa(r.Loaded))
	builder.WriteRune('/')
	builder.WriteString(strconv.Itoa(r.Filtered))
	builder.WriteString("</code>")

	return builder.String()
}
