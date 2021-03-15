package loader

import (
	"gfeed/utils"
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

type Element interface {
	Attr(attribute string) string
	ChildAttr(selector, attribute string) string
	ChildText(selector string) string
}

func (option PathFinder) findAttribute(e Element) string {
	val := option.findAttributeRaw(e)

	if option.ParseStrategy == ParserStrategyStyle {
		return parseStyle((val))
	}

	return val
}

func (option PathFinder) findAttributeRaw(e Element) string {
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

func (category PathFinderCategory) isAllowed(categories []string) bool {
	for _, v := range categories {
		_, ok := utils.FindStr(category.Alloweds, strings.ToLower(v))

		if ok {
			return true
		}
	}

	return false
}

func (category PathFinderCategory) findCategories(e interface{}) []string {
	if len(category.Path) == 0 {
		return []string{}
	}

	xml, ok := e.(*colly.XMLElement)

	if ok {
		return category.findCategories_onXML(xml)
	}

	html := e.(*colly.HTMLElement)

	return category.findCategories_onHTML(html)

}

func (category PathFinderCategory) findCategories_onXML(e *colly.XMLElement) (cats []string) {
	return e.ChildTexts(category.Path)
}

func (category PathFinderCategory) findCategories_onHTML(e *colly.HTMLElement) (cats []string) {
	e.ForEach(category.Path, func(i int, h *colly.HTMLElement) {
		if len(category.Attribute) == 0 {
			cats = append(cats, h.Text)
			return
		}

		cats = append(cats, h.Attr(category.Attribute))
	})

	return cats
}
