package loader

import (
	log "gfeed/utils/logger"
)

var baseLogger log.Logger

func init() {
	baseLogger = log.New("scrapper:loader")
}
