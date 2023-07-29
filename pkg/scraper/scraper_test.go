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

func (s *FindEntriesTestSuite) TestExample01() {
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

func (s *FindEntriesTestSuite) TestExample02() {
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

func TestFindEntriesSuite(t *testing.T) {
	suite.Run(t, new(FindEntriesTestSuite))
}
