package testdata

import (
	"embed"
	"net/http"
	"strings"

	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"gopkg.in/yaml.v3"
)

//go:embed *.html
//go:embed *.xml
//go:embed *.json
var files embed.FS

func ParseSource(baseURL, input string) (scraper.SourceDefinition, error) {
	var source scraper.SourceDefinition

	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\t", "  ")
	input = strings.TrimSpace(input)

	err := yaml.Unmarshal([]byte(input), &source)

	if err != nil {
		return source, err
	}

	source.BaseURL = baseURL

	return source, nil
}

func FileHandler() http.Handler {
	return http.FileServer(http.FS(files))
}
