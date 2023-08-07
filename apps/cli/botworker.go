package main

import (
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/apps/cli/actions"
)

func botCMD() *cli.Command {
	worker := &cli.Command{
		Name:        "worker",
		Description: `Start bot worker`,
		Action: func(c *cli.Context) error {
			return actions.BotWorker(c.Context)
		},
	}

	return &cli.Command{
		Name:        "bot",
		Subcommands: []*cli.Command{worker},
	}
}
