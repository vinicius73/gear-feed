package scraper

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

var reURL = regexp.MustCompile(`\((.*?)\)`)

type PathParserStrategy = string

const (
	ParserStrategyNone  PathParserStrategy = ""
	ParserStrategyStyle PathParserStrategy = "style"
)

const (
	XML  = "XML"
	HTML = "HTML"
)

type SourceDefinition struct {
	Name       string           `yaml:"name"`
	Enabled    bool             `yaml:"enabled"`
	BaseURL    string           `yaml:"base_url"`
	Path       string           `yaml:"path"`
	Limit      int8             `yaml:"limit"`
	Parser     string           `yaml:"parser"`
	Attributes AttributesFinder `yaml:"attributes"`
}

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
	EntrySelector string             `yaml:"entry_selector"`
	Category      PathFinderCategory `yaml:"category"`
	Link          PathFinder         `yaml:"link"`
	Title         PathFinder         `yaml:"title"`
	Image         PathFinder         `yaml:"image"`
}

func parseStyle(style string) string {
	result := reURL.FindString(style)

	result = strings.TrimLeft(result, "(")
	result = strings.TrimRight(result, ")")

	if strings.HasPrefix(result, "//") {
		return "https:" + result
	}

	return result
}

type Element interface {
	Attr(attribute string) string
	ChildAttr(selector, attribute string) string
	ChildText(selector string) string
}

func (d SourceDefinition) visitURL() string {
	return d.BaseURL + d.Path
}

func (d SourceDefinition) buildEntry(title, link, image string, categories []string) Entry {
	return Entry{
		Type:       d.Name,
		Title:      title,
		Categories: categories,
		Link:       d.absouteURL(link),
		Image:      d.absouteURL(image),
	}
}

func (d SourceDefinition) absouteURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}

	if strings.HasPrefix(path, "//") {
		return "https:" + path
	}

	return d.BaseURL + path
}

func (option PathFinder) findAttribute(e Element) string {
	val := option.findAttributeRaw(e)

	if option.ParseStrategy == ParserStrategyStyle {
		val = parseStyle(val)
	}

	return val
}

func (option PathFinder) findAttributeRaw(el Element) string {
	if len(option.Path) == 0 {
		if len(option.Attribute) > 0 {
			return el.Attr(option.Attribute)
		}

		return ""
	}

	if len(option.Attribute) > 0 {
		return el.ChildAttr(option.Path, option.Attribute)
	}

	return el.ChildText(option.Path)
}

func (category PathFinderCategory) isAllowed(categories []string) bool {
	if len(category.Alloweds) == 0 {
		return true
	}

	return support.ContainsSome(category.Alloweds, categories)
}

func (category PathFinderCategory) findCategories(element interface{}) []string {
	if len(category.Path) == 0 {
		return []string{}
	}

	xml, ok := element.(*colly.XMLElement)

	if ok {
		return category.findCategoriesOnXML(xml)
	}

	html, ok := element.(*colly.HTMLElement)

	if ok {
		return category.findCategoriesOnHTML(html)
	}

	return []string{}
}

func (category PathFinderCategory) findCategoriesOnXML(element *colly.XMLElement) []string {
	return element.ChildTexts(category.Path)
}

func (category PathFinderCategory) findCategoriesOnHTML(element *colly.HTMLElement) []string {
	cats := []string{}

	element.ForEach(category.Path, func(i int, child *colly.HTMLElement) {
		if len(category.Attribute) == 0 {
			cats = append(cats, child.Text)

			return
		}

		cats = append(cats, child.Attr(category.Attribute))
	})

	return cats
}
