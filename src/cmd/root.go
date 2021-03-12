package cmd

import (
	"fmt"
	"gfeed/bot"
	"gfeed/utils"
	"os"

	"github.com/spf13/cobra"
)

var token string
var channel string
var user string
var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "gfeed",
	Short: "Gamer Feed Project",
	Run: func(cmd *cobra.Command, args []string) {
		bot.SendNews(getBotConfig())
	},
}

func init() {
	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&user, "user", "u", utils.GetEnv("TELEGRAM_USER", ""), "Telegram User")
	flags.StringVarP(&token, "token", "t", "", "Telegram Token (required)")
	flags.StringVarP(&channel, "channel", "c", "@GamerFeed", "Telegram Channel")
	flags.BoolVarP(&dryRun, "dry", "", false, "Just try to run")

	rootCmd.MarkFlagRequired("token")
}

// Execute the process
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
