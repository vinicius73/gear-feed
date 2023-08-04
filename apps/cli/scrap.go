package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/pkg/configurations"
	"github.com/vinicius73/gamer-feed/pkg/linkloader"
	"github.com/vinicius73/gamer-feed/pkg/scraper"
	"github.com/vinicius73/gamer-feed/pkg/sender"
	"github.com/vinicius73/gamer-feed/pkg/telegram"
	"github.com/vinicius73/gamer-feed/sources"
)

func loadEntries(ctx context.Context, only []string) ([]scraper.Entry, error) {
	list, err := sources.LoadDefinitions(ctx, sources.LoadOptions{
		Only: only,
	})
	if err != nil {
		return []scraper.Entry{}, err
	}

	collections, err := linkloader.FromSources(ctx, linkloader.LoadOptions{
		Workers: (len(list) + 1) / 2,
		Sources: list,
	})
	if err != nil {
		return []scraper.Entry{}, err
	}

	return collections.Shuffle(), nil
}

func buildSender(ctx context.Context, chats []int64) (sender.Serder, error) {
	config := configurations.Ctx(ctx)

	bot, err := telegram.NewBot(config.Telegram)
	if err != nil {
		return nil, err
	}

	return sender.NewTelegramSerder(bot, chats), nil
}

func scrapCMD() *cli.Command {
	load := &cli.Command{
		Name:        "load",
		Description: `Load scrap data from sources.`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Usage:   "Limit the number of entries to send",
				Aliases: []string{"l"},
				Value:   10,
				Action: func(c *cli.Context, val int) error {
					if val < 1 {
						return cli.Exit("Limit must be greater than 1", 1)
					}

					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:    "only",
				Usage:   "Load only the specified loaders",
				Aliases: []string{"o"},
			},
			&cli.Int64Flag{
				Name:     "to",
				Usage:    "Send the loaded data to the specified channel",
				Required: true,
			},
		},
		Action: func(cmd *cli.Context) error {
			logger := zerolog.Ctx(cmd.Context)

			only := cmd.StringSlice("only")
			limit := cmd.Int("limit")
			to := cmd.Int64("to")

			entries, err := loadEntries(cmd.Context, only)
			if err != nil {
				return err
			}

			logger.Info().Msgf("Found %d entries", len(entries))

			if limit > 0 && len(entries) > limit {
				logger.Warn().Msgf("Limiting to %d entries", limit)
				entries = entries[:limit]
			}

			botSender, err := buildSender(cmd.Context, []int64{to})
			if err != nil {
				return err
			}

			sendables := make([]sender.Sendable, len(entries))

			for index, entry := range entries {
				sendables[index] = sender.NewScrapEntry(entry)
			}

			botSender.SendCollection(cmd.Context, sendables)

			return nil
		},
	}

	return &cli.Command{
		Name:        "scrap",
		Subcommands: []*cli.Command{load},
	}
}
