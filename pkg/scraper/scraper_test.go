package scraper_test

import (
	"context"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/scraper/testdata"
	"gopkg.in/yaml.v3"
)

type FindEntriesTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *FindEntriesTestSuite) SetupTest() {
	s.server = httptest.NewServer(testdata.FileHandler())
}

func (s *FindEntriesTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *FindEntriesTestSuite) parseSource(input string) scraper.SourceDefinition {
	var source scraper.SourceDefinition

	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\t", "  ")
	input = strings.TrimSpace(input)

	err := yaml.Unmarshal([]byte(input), &source)

	s.NoError(err)

	source.BaseURL = s.server.URL

	return source
}

func (s *FindEntriesTestSuite) TestExample01Simple() {
	source := s.parseSource(`
name: test
enabled: true
path: /example_01.html
attributes:
	entry_selector: "#news > article"
	link:
		path: "h2 a"
		attribute: "href"
	title:
		path: "h2 a"
	image:
		path: "figure img"
		attribute: "src"
	`)

	entries, err := scraper.FindEntries(context.TODO(), source)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 3, len(entries))

	for index, entry := range entries {
		assert.Equal(s.T(), "Good news "+strconv.Itoa(index+1), entry.Title)
		assert.Equal(s.T(), "http://foo.com/news/good-"+strconv.Itoa(index+1), entry.Link)
		assert.Equal(s.T(), "http://bar.bang/foo.jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample02BaseURL() {
	source := s.parseSource(`
name: test_01
enabled: true
path: /example_02.html
attributes:
	entry_selector: "#posts > article"
	link:
		path: "a"
		attribute: "href"
	title:
		path: "h2"
	image:
		path: "img"
		attribute: "data-src"
	`)

	entries, err := scraper.FindEntries(context.TODO(), source)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 2, len(entries))

	baseURL := source.BaseURL

	for index, entry := range entries {
		num := strconv.Itoa(index + 1)
		assert.Equal(s.T(), "Hot News "+num, entry.Title)
		assert.Equal(s.T(), baseURL+"/hot-"+num+".htm", entry.Link)
		assert.Equal(s.T(), baseURL+"/bang-"+num+".png", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample03XML() {
	source := s.parseSource(`
name: test_xml
path: /example_03.xml
limit: 2
parser: XML
attributes:
	entry_selector: //channel[1]/item
	link:
		path: /link
	title:
		path: /title
	image:
		path: enclosure
		attribute: url
	
	`)

	entries, err := scraper.FindEntries(context.TODO(), source)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 3, len(entries))

	for index, entry := range entries {
		num := strconv.Itoa(index + 1)
		assert.Equal(s.T(), "XML in 200"+num, entry.Title)
		assert.Equal(s.T(), "https://xmlsite.net/news-"+num+".html", entry.Link)
		assert.Equal(s.T(), "https://xmlsite.net/news-"+num+".jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample04CategoriesFilter() {
	source := s.parseSource(`
name: test_04
enabled: true
path: /example_04.html
attributes:
	entry_selector: "#last-news-games > article"
	link:
		path: "h2 a"
		attribute: "href"
	title:
		path: "h2 a"
	category:
		path_finder:
			path: ul.tags > li > a
		allows:
			- "d"
	image:
		path: "figure img"
		attribute: "src"
	`)

	entries, err := scraper.FindEntries(context.TODO(), source)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 2, len(entries))

	for index, num := range []string{"2", "3"} {
		entry := entries[index]

		assert.Equal(s.T(), "Good news "+num, entry.Title)
		assert.Equal(s.T(), "http://foo.com/games/good-"+num, entry.Link)
		assert.Equal(s.T(), "https://super.site/foo.jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample05XMLCategories() {
	source := s.parseSource(`
name: test_xml_categories
path: /example_05.xml
limit: 2
parser: XML
attributes:
	entry_selector: //channel[1]/item
	category:
		path_finder:
			path: /category
		allows:
			- "a1"
			- "h7"
	link:
		path: /link
	title:
		path: /title
	image:
		path: enclosure
		attribute: url
	
	`)

	entries, err := scraper.FindEntries(context.TODO(), source)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 4, len(entries))

	for index, num := range []string{"1", "3", "4", "6"} {
		entry := entries[index]
		assert.Equal(s.T(), "XML Title in 200"+num, entry.Title)
		assert.Equal(s.T(), "https://xmlsite.net/news-"+num+".html", entry.Link)
		assert.Equal(s.T(), "https://xmlsite.net/news-"+num+".jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample06CategoriesFromAttributesFilter() {
	source := s.parseSource(`
name: test_06_categories_from_attributes
enabled: true
path: /example_06.html
attributes:
	entry_selector: "#last-news-games > article"
	link:
		path: "h2 a"
		attribute: "href"
	title:
		path: "h2 a"
	category:
		path_finder:
			path: ul.tags > li
			attribute: "data-category"
		allows:
			- "cat-b"
			- "cat-c"
	image:
		path: "figure img"
		attribute: "src"
	`)

	entries, err := scraper.FindEntries(context.TODO(), source)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 2, len(entries))

	for index, num := range []string{"1", "5"} {
		entry := entries[index]

		assert.Equal(s.T(), "Game news "+num, entry.Title)
		assert.Equal(s.T(), "http://foo.com/games/good-"+num, entry.Link)
		assert.Equal(s.T(), "https://super.site/baz.jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample07CustomAttributeParser() {
	source := s.parseSource(`
name: test
enabled: true
path: /example_07.html
attributes:
	entry_selector: "#news > article"
	link:
		path: "h2 a"
		attribute: "href"
	title:
		path: "h2 a"
		attribute: "alt"
	image:
		path: "figure"
		attribute: "style"
		parse_strategy: "style"
	`)

	entries, err := scraper.FindEntries(context.TODO(), source)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 3, len(entries))

	for index, entry := range entries {
		num := strconv.Itoa(index + 1)
		title := "G@M3R news " + num + " Não há quem goste de dor, que a procure e a queira ter,"
		assert.Equal(s.T(), title, entry.Title)
		assert.Equal(s.T(), "http://foo.com/news/good-"+num, entry.Link)
		assert.Equal(s.T(), "https://cdn.net/images/news-"+num+".png", entry.Image)
	}
}

func TestFindEntriesSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FindEntriesTestSuite))
}
