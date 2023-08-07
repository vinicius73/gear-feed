package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	zero "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/gamer-feed/pkg"
	"github.com/vinicius73/gamer-feed/pkg/configurations"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

var fileLog *os.File

func main() {
	var err error

	defer func() {
		if fileLog != nil {
			fileLog.Close()
		}

		if err != nil {
			zero.Fatal().Err(err).Msg("Fail run application")
			os.Exit(1)
		}
	}()

	mainCtx, cancel := support.WithKillSignal(context.Background())

	defer cancel()

	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "gfeed",
		Usage:                "Gamer Feed Bot CLI",
		Version:              pkg.Version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from",
				DefaultText: fmt.Sprintf("%s/gfeed.yml", support.GetBinDirPath()),
			},
			&cli.StringFlag{
				Name:        "level",
				Aliases:     []string{"l"},
				Usage:       "define log level",
				DefaultText: "info",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enables debug level",
			},
			&cli.StringFlag{
				Name:  "log-file",
				Value: "",
				Usage: "store logs in a file",
			},
		},
		Commands: []*cli.Command{
			sourcesCMD(),
			scrapCMD(),
			botCMD(),
		},
		Before: beforeRun,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	// catch on defer
	//nolint:errcheck
	err = app.RunContext(mainCtx, os.Args)
}

func beforeRun(cmd *cli.Context) error {
	debug := cmd.Bool("debug")

	appConfig, err := configurations.Load(cmd.String("config"))
	if err != nil {
		e, ok := err.(apperrors.BusinessError) //nolint:errorlint
		if ok && e.ErrorCode == configurations.ConfigFileWasCreated.ErrorCode {
			zero.Warn().Msg(err.Error())
		} else {
			zero.Fatal().Err(err).Msg("Fail to load config")

			return err
		}
	}

	logFile := cmd.String("log-file")

	if logFile != "" {
		fileLog, err = support.LoggerToFile(logFile)
		if err != nil {
			return err
		}
	}

	logLevel := cmd.String("level")

	if debug {
		logLevel = "debug"
	}

	if logLevel != "" {
		appConfig.Logger.Level = logLevel
	}

	if appConfig.Logger.Debug(logLevel) {
		appConfig.Debug = true
	}

	cmd.Context = appConfig.WithContext(cmd.Context)

	support.SetupLogger(appConfig.Logger.Level, appConfig.Logger.Format, appConfig.Tags())

	log := support.Logger("", appConfig.Tags())

	// inject logger into context
	ctx := log.WithContext(cmd.Context)
	ctx = log.WithContext(ctx)

	cmd.Context = ctx

	log.Debug().Msg("Application started")

	return nil
}
