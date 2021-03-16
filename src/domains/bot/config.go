package bot

import "gfeed/domains"

// Config of bot
type Config struct {
	Token   string
	Channel string
	User    string
	Info    domains.ProcessInfo
	DryRun  bool
}
