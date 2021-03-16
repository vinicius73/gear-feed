package cmd

import (
	"fmt"
	"gfeed/domains"
	"gfeed/utils/logger"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version gamer feed cli",
	Run: func(cmd *cobra.Command, args []string) {
		info := domains.Info()
		fmt.Printf(
			"GemerFeed by @vinicius73 v%s -- %s / %s", info.Version, info.Commit, info.BuildDate,
		)
	},
}

func versionHook(cmd *cobra.Command, args []string) {
	info := domains.Info()

	logger.Global().
		Info().
		Str("version", info.Version).
		Str("commit", info.Commit).
		Str("buildDate", info.BuildDate).
		Msg("Running Gamer Feed")
}
