package logger

import (
	"gfeed/utils"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger = zerolog.Logger

func init() {
	log.Logger = getLogFormatter(log.Logger)
	zerolog.SetGlobalLevel(getLogLevel())
}

func New(context string) zerolog.Logger {
	return log.
		With().
		Str("context", context).
		Logger()
}

func getLogLevel() zerolog.Level {
	level := strings.ToLower(utils.GetEnv("LOG_LEVEL", "info"))

	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

func getLogFormatter(l zerolog.Logger) zerolog.Logger {
	level := strings.ToLower(utils.GetEnv("LOG_FORMAT", "text"))

	switch level {
	case "json":
		return l
	default:
		return l.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
