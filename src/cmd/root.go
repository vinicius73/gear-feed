package cmd

import (
	"fmt"
	"gfeed/bot"
	"os"

	"github.com/spf13/cobra"
)

var token string
var channel string

var rootCmd = &cobra.Command{
	Use:   "gfeed",
	Short: "Gamer Feed Project",
	Run: func(cmd *cobra.Command, args []string) {
		bot.SendNews(bot.Config{
			Token:   token,
			Channel: channel,
		})
	},
}

func init() {
	flags := rootCmd.Flags()

	flags.StringVarP(&token, "token", "t", "", "Telegram Token (required)")
	flags.StringVarP(&channel, "channel", "c", "@GamerFeed", "Telegram Channel")

	rootCmd.MarkFlagRequired("token")
}

// Execute the process
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
