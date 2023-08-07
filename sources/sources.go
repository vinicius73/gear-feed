package sources

import (
	"context"
	"embed"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"gopkg.in/yaml.v3"
)

//go:embed *.yml
var files embed.FS

type LoadOptions struct {
	Only []string
}

func LoadDefinitions(ctx context.Context, options LoadOptions) ([]scraper.SourceDefinition, error) {
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

		switch {
		case len(options.Only) > 0:
			if support.Contains(options.Only, def.Name) {
				definitions = append(definitions, def)
			} else {
				logger.Warn().
					Msgf("Loader %s is not in the list of loaders to be loaded", def.Name)
			}
		case def.Enabled:
			definitions = append(definitions, def)
		default:
			logger.Warn().
				Msgf("Loader %s is disabled", def.Name)
		}
	}

	return definitions, nil
}
