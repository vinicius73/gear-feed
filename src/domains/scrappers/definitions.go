package scrappers

import (
	"embed"
	"gfeed/domains/news"
	"gfeed/domains/scrappers/loader"

	"gopkg.in/yaml.v2"
)

type loaderFn = func() []news.Entry

//go:embed definitions/*.yml
var files embed.FS

func loadDefinitions() (definitions []loader.Definitions, err error) {
	dir, err := files.ReadDir("definitions")

	if err != nil {
		return definitions, err
	}

	for _, entry := range dir {
		file, _ := files.ReadFile("definitions/" + entry.Name())

		var def loader.Definitions

		err = yaml.Unmarshal(file, &def)

		if err != nil {
			return definitions, err
		}

		if def.Enabled {
			definitions = append(definitions, def)
		} else {
			logger.Warn().
				Msgf("Loasder %s is disabled", def.Name)
		}

	}

	return definitions, nil
}
