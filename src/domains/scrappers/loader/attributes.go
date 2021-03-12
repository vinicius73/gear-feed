package loader

import (
	"strings"

	"github.com/gocolly/colly"
)

type PathFinder struct {
	Path          string             `yaml:"path"`
	Attribute     string             `yaml:"attribute"`
	ParseStrategy PathParserStrategy `yaml:"parse_strategy"`
}

type PathFinderCategory struct {
	PathFinder `yaml:"path_finder"`
	Alloweds   []string `yaml:"allows"`
}

type AttributesFinder struct {
	Wrapper  string             `yaml:"wrapper"`
	Category PathFinderCategory `yaml:"category"`
	Link     PathFinder         `yaml:"link"`
	Title    PathFinder         `yaml:"title"`
	Image    PathFinder         `yaml:"image"`
}

func (option PathFinder) findAttribute(e *colly.HTMLElement) string {
	val := option.findAttributeRaw(e)

	if option.ParseStrategy == ParserStrategyStyle {
		return parseStyle((val))
	}

	return val
}

func (option PathFinder) findAttributeRaw(e *colly.HTMLElement) string {
	if len(option.Path) == 0 {
		if len(option.Attribute) > 0 {
			return e.Attr(option.Attribute)
		}

		return ""
	}

	if len(option.Attribute) > 0 {
		return e.ChildAttr(option.Path, option.Attribute)
	}

	return e.ChildText(option.Path)
}

func (category PathFinderCategory) isAllowed(cat string) bool {
	c := strings.ToLower(cat)

	for _, v := range category.Alloweds {
		if c == v {
			return true
		}
	}

	return false
}
