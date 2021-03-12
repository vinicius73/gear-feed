package cmd

import (
	"fmt"
	"gfeed/utils/logger"

	"github.com/spf13/cobra"
)

type ProcessInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version gamer feed cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			"GemerFeed by @vinicius73 v%s -- %s / %s", info.Version, info.Commit, info.BuildDate,
		)
	},
}

func versionHook(cmd *cobra.Command, args []string) {
	logger.Global().
		Info().
		Str("version", info.Version).
		Str("commit", info.Commit).
		Str("buildDate", info.BuildDate).
		Msg("Running Gamer Feed")
}
