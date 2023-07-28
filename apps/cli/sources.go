package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/pkg/support"
	layout "github.com/vinicius73/gamer-feed/pkg/tui/sources"
	"github.com/vinicius73/gamer-feed/sources"
)

func sourcesCMD() *cli.Command {
	list := &cli.Command{
		Name: "list",
		Action: func(cmd *cli.Context) error {
			list, err := sources.LoadDefinitions(cmd.Context)
			if err != nil {
				return err
			}

			logFile := cmd.String("log-file")

			if logFile != "" {
				f, err := support.LoggerToFile(logFile)
				if err != nil {
					return err
				}

				defer f.Close()
			}

			logger := support.Logger("tui", nil)

			ctx := logger.WithContext(cmd.Context)

			p := tea.NewProgram(
				layout.NewModel(ctx, list),
				tea.WithAltScreen(),
				tea.WithContext(ctx),
			)
			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return &cli.Command{
		Name:        "sources",
		Description: "Interact with sources",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-file",
				Value: "",
				Usage: "Store logs in a file",
			},
		},
		Subcommands: []*cli.Command{list},
	}
}
