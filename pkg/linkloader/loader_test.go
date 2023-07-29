package linkloader_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/scraper/testdata"
)

const sourceDemo01 = `
name: Demo 01
enabled: true
path: /demo_01.html
attributes:
	entry_selector: "body ul li"
	link:
		path: "a"
		attribute: "href"
	title:
		path: "a"
	image:
		path: "img"
		attribute: "src"
`

const sourceDemo02 = `
name: Demo 02
enabled: true
path: /demo_02.html
attributes:
	entry_selector: "body ul li"
	link:
		path: "a"
		attribute: "href"
	title:
		path: "a"
	image:
		path: "img"
		attribute: "src"
`

func TestFromSources(t *testing.T) {
	server := httptest.NewServer(testdata.FileHandler())

	defer server.Close()

	demo01, err := testdata.ParseSource(server.URL, sourceDemo01)

	assert.NoError(t, err)

	demo02, err := testdata.ParseSource(server.URL, sourceDemo02)

	assert.NoError(t, err)

	sources := []scraper.SourceDefinition{demo01, demo02}

	entries, err := linkloader.FromSources(context.TODO(), linkloader.LoadOptions{
		Sources: sources,
		Workers: 2,
	})

	assert.NoError(t, err)

	assert.Equal(t, 2, len(entries))
	assert.Equal(t, 6, len(entries.Entries()))
}
