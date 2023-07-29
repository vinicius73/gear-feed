package main

import (
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/sources"
)

func scrapCMD() *cli.Command {
	load := &cli.Command{
		Name:        "load",
		Description: `Load scrap data from sources.`,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "only",
				Usage:   "Load only the specified loaders",
				Aliases: []string{"o"},
			},
		},
		Action: func(cmd *cli.Context) error {
			logger := zerolog.Ctx(cmd.Context)

			only := cmd.StringSlice("only")

			list, err := sources.LoadDefinitions(cmd.Context, sources.LoadOptions{
				Only: only,
			})
			if err != nil {
				return err
			}

			collections, err := linkloader.FromSources(cmd.Context, linkloader.LoadOptions{
				Workers: (len(list) + 1) / 2,
				Sources: list,
			})
			if err != nil {
				return err
			}

			entries := collections.Entries()

			logger.Info().Msgf("Found %d entries", len(entries))

			// for _, entry := range entries {
			// 	logger.Info().Msgf("Entry: %s", entry.Link)
			// }

			return nil
		},
	}

	return &cli.Command{
		Name:        "scrap",
		Subcommands: []*cli.Command{load},
		Before: func(cmd *cli.Context) error {
			logger := support.Logger("scrap", nil)

			cmd.Context = logger.WithContext(cmd.Context)

			return nil
		},
	}
}
