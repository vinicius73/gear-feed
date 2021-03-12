package loader

import (
	"regexp"
	"strings"
)

type PathParserStrategy = int

const (
	ParserStrategyNone PathParserStrategy = iota
	ParserStrategyStyle
)

var reURL = regexp.MustCompile(`\((.*?)\)`)

func parseStyle(style string) string {
	result := reURL.FindString(style)

	result = strings.TrimLeft(result, "(")
	result = strings.TrimRight(result, ")")

	if strings.HasPrefix(result, "//") {
		return "https:" + result
	}

	return result

}
