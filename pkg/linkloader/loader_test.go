package linkloader_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/gear-feed/pkg/linkloader"
	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/scraper"
	"github.com/vinicius73/gear-feed/pkg/scraper/testdata"
)

const sourceDemo01 = `
name: Demo 01
enabled: true
paths:
  - /demo_01.html
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
limit: 4
paths:
  - /demo_01.html
  - /demo_02.html
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

// TestFromSources tests the FromSources function of the linkloader package.
// It sets up an HTTP test server to serve test data, parses two source definitions,
// and verifies that the FromSources function correctly loads entries from these sources.
// The test ensures that no errors occur during the loading process and that the expected
// number of entries and entries' length are returned.
func TestFromSources(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(testdata.FileHandler())

	defer server.Close()

	demo01, err := testdata.ParseSource(server.URL, sourceDemo01)

	assert.NoError(t, err)

	demo02, err := testdata.ParseSource(server.URL, sourceDemo02)

	assert.NoError(t, err)

	sources := []scraper.SourceDefinition{demo01, demo02}

	entries, err := linkloader.FromSources[model.Entry](context.TODO(), linkloader.LoadOptions{
		Sources: sources,
		Workers: 2,
	})

	assert.NoError(t, err)

	assert.Len(t, entries, 2)
	// The first source has 3 entries
	// The second source has 6 entries, but the limit is set to 4
	assert.Len(t, entries.Entries(), 7)
}

func TestLoadEntries(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(testdata.FileHandler())

	defer server.Close()

	demo01, err := testdata.ParseSource(server.URL, sourceDemo01)

	assert.NoError(t, err)

	demo02, err := testdata.ParseSource(server.URL, sourceDemo02)

	assert.NoError(t, err)

	sources := []scraper.SourceDefinition{demo01, demo02}

	entries, err := linkloader.LoadEntries[model.Entry](context.TODO(), linkloader.LoadOptions{
		Workers: 0,
		Sources: sources,
	})

	assert.NoError(t, err)

	assert.Len(t, entries, 7)
}
