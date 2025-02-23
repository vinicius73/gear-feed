package main

import (
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gear-feed/apps/cli/actions"
)

func dbCMD() *cli.Command {
	cleanup := &cli.Command{
		Name:        "cleanup",
		Description: `Cleanup the database.`,
		Action: func(cmd *cli.Context) error {
			return actions.Cleanup(cmd.Context)
		},
	}

	return &cli.Command{
		Name:        "db",
		Description: `Database related commands.`,
		Subcommands: []*cli.Command{
			cleanup,
		},
	}
}
