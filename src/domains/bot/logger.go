package bot

import (
	log "gfeed/utils/logger"
)

var logger log.Logger

func init() {
	logger = log.New("bot")
}
