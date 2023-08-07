package main

import (
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/apps/cli/actions"
)

func scrapCMD() *cli.Command {
	load := &cli.Command{
		Name:        "load",
		Description: `Load scrap data from sources.`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Usage:   "Limit the number of entries to send",
				Aliases: []string{"l"},
				Value:   10, //nolint:gomnd // default value
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
			only := cmd.StringSlice("only")
			limit := cmd.Int("limit")
			to := cmd.Int64("to")

			return actions.Load(cmd.Context, actions.LoadOptions{
				Only:  only,
				To:    to,
				Limit: limit,
			})
		},
	}

	return &cli.Command{
		Name:        "scrap",
		Subcommands: []*cli.Command{load},
	}
}
