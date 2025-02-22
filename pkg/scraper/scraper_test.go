package scraper_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/scraper/testdata"
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
	source, err := testdata.ParseSource(s.server.URL, input)

	s.NoError(err)

	return source
}

func (s *FindEntriesTestSuite) TestExample01Simple() {
	source := s.parseSource(`
name: test
enabled: true
paths:
 - /example_01.html
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

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)

	s.NoError(err)

	s.Len(entries, 3)

	for index, entry := range entries {
		s.Equal("Good news "+strconv.Itoa(index+1), entry.Title)
		s.Equal("http://foo.com/news/good-"+strconv.Itoa(index+1), entry.URL)
		s.Equal("http://bar.bang/foo.jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample02BaseURL() {
	source := s.parseSource(`
name: test_01
enabled: true
base_url: "http://json.com"
paths:
  - /example_02.html
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

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)

	s.NoError(err)

	s.Len(entries, 2)

	baseURL := source.BaseURL

	for index, entry := range entries {
		num := strconv.Itoa(index + 1)
		s.Equal("Hot News "+num, entry.Title)
		s.Equal(baseURL+"/hot-"+num+".htm", entry.URL)
		s.Equal(baseURL+"/bang-"+num+".png", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample03XML() {
	source := s.parseSource(`
name: test_xml
paths:
  - /example_03.xml
limit: 0
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

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)

	s.NoError(err)

	s.Len(entries, 3)

	for index, entry := range entries {
		num := strconv.Itoa(index + 1)
		s.Equal("XML in 200"+num, entry.Title)
		s.Equal("https://xmlsite.net/news-"+num+".html", entry.URL)
		s.Equal("https://xmlsite.net/news-"+num+".jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample04CategoriesFilter() {
	source := s.parseSource(`
name: test_04
enabled: true
paths:
  - /example_04.html
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

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)

	s.NoError(err)

	s.Len(entries, 2)

	for index, num := range []string{"2", "3"} {
		entry := entries[index]

		s.Equal("Good news "+num, entry.Title)
		s.Equal("http://foo.com/games/good-"+num, entry.URL)
		s.Equal("https://super.site/foo.jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample05XMLCategories() {
	source := s.parseSource(`
name: test_xml_categories
paths:
  - /example_05.xml
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

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)

	s.NoError(err)

	s.Len(entries, 2)

	for index, num := range []string{"1", "3"} {
		entry := entries[index]
		s.Equal("XML Title in 200"+num, entry.Title)
		s.Equal("https://xmlsite.net/news-"+num+".html", entry.URL)
		s.Equal("https://xmlsite.net/news-"+num+".jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample06CategoriesFromAttributesFilter() {
	source := s.parseSource(`
name: test_06_categories_from_attributes
enabled: true
paths:
  - /example_06.html
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

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)

	s.NoError(err)

	s.Len(entries, 2)

	for index, num := range []string{"1", "5"} {
		entry := entries[index]

		s.Equal("Game news "+num, entry.Title)
		s.Equal("http://foo.com/games/good-"+num, entry.URL)
		s.Equal("https://super.site/baz.jpg", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExample07CustomAttributeParser() {
	source := s.parseSource(`
name: test
enabled: true
paths:
  - /example_07.html
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

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)

	s.NoError(err)

	s.Len(entries, 3)

	for index, entry := range entries {
		num := strconv.Itoa(index + 1)
		title := "G@M3R news " + num + " Não há quem goste de dor, que a procure e a queira ter, simplesmente porque é dor..."
		s.Equal(title, entry.Title)
		s.Equal("http://foo.com/news/good-"+num, entry.URL)
		s.Equal("https://cdn.net/images/news-"+num+".png", entry.Image)
	}
}

func (s *FindEntriesTestSuite) TestExampleJSON() {
	source := s.parseSource(`
  name: JSONSOURCE
paths:
  - /example.json

limit: 3
enabled: true
parser: JSON
attributes:
  entry_selector: "stories"
  link:
    path: slug
  title:
    path: content.lead
  image:
    path: content.thumbnail.filename
`)

	entries, err := scraper.FindEntries[model.Entry](context.TODO(), source)
	s.NoError(err)
	s.Len(entries, 3)

	baseURL := source.BaseURL

	for index, entry := range entries {
		num := strconv.Itoa(index + 1)

		s.Equal("A new Entry 00"+num, entry.Title)
		s.Equal(fmt.Sprintf("%s/latest/2023/9/24/new-entries-00%v", baseURL, num), entry.URL)
		s.Equal("https://foo.json/image-00"+num+".jpg", entry.Image)
	}
}

func TestFindEntriesSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FindEntriesTestSuite))
}
