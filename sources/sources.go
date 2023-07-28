package sources

import (
	"context"
	"embed"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"gopkg.in/yaml.v2"
)

//go:embed *.yml
var files embed.FS

func LoadDefinitions(ctx context.Context) ([]scraper.SourceDefinition, error) {
	definitions := []scraper.SourceDefinition{}

	dir, err := files.ReadDir(".")
	if err != nil {
		return definitions, err
	}

	logger := zerolog.Ctx(ctx)

	for _, entry := range dir {
		file, _ := files.ReadFile(entry.Name())

		var def scraper.SourceDefinition

		err = yaml.Unmarshal(file, &def)

		if err != nil {
			return definitions, err
		}

		if def.Enabled {
			definitions = append(definitions, def)
		} else {
			logger.Warn().
				Msgf("Loader %s is disabled", def.Name)
		}
	}

	return definitions, nil
}
