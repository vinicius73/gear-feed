package sources

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
	"gopkg.in/yaml.v3"
)

var ErrFailToLoad = apperrors.System(nil, "fail to load sources", "SOURCES:FAIL_TO_LOAD")

var ymlExtensions = []string{".yml", ".yaml"}

type LoadOptions struct {
	Paths []string `fig:"paths" yaml:"paths"`
	Only  []string `fig:"only"  yaml:"only"`
}

func Load(ctx context.Context, opt LoadOptions) (Collection, error) {
	definitions := []scraper.SourceDefinition{}

	for _, path := range opt.Paths {
		defs, err := loadPath(ctx, path, opt.Only)
		if err != nil {
			return definitions, ErrFailToLoad.Wrap(err)
		}

		definitions = append(definitions, defs...)
	}

	return definitions, nil
}

func loadPath(ctx context.Context, path string, only []string) (Collection, error) {
	definitions := []scraper.SourceDefinition{}

	if !filepath.IsAbs(path) {
		pwd, err := os.Getwd()
		if err != nil {
			return definitions, ErrFailToLoad.Wrap(err)
		}

		path = filepath.Join(pwd, path)
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return definitions, err
	}

	logger := zerolog.Ctx(ctx).With().Str("path", path).Logger()

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		base := filepath.Ext(name)

		if !support.Contains(ymlExtensions, strings.ToLower(base)) {
			logger.Warn().
				Msgf("File %s is not a yaml file", name)

			continue
		}

		def, use, err := readFile(filepath.Join(path, name), only)
		if err != nil {
			return definitions, ErrFailToLoad.Wrap(err)
		}

		if use {
			definitions = append(definitions, def)
		} else {
			logger.Debug().
				Msgf("Ignoring %s", def.Name)
		}
	}

	return definitions, nil
}

func readFile(fileName string, only []string) (scraper.SourceDefinition, bool, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return scraper.SourceDefinition{}, false, err
	}

	var def scraper.SourceDefinition

	err = yaml.Unmarshal(file, &def)

	if err != nil {
		return scraper.SourceDefinition{}, false, err
	}

	if len(only) > 0 {
		if support.Contains(only, def.Name) {
			return def, true, nil
		}

		return def, false, nil
	}

	if def.Enabled {
		return def, true, nil
	}

	return def, false, nil
}
