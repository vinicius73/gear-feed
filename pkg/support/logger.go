package support

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vinicius73/gear-feed/pkg"
)

//nolint:gochecknoinits
func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldInteger = true
}

func SetupLogger(level, format string, tags map[string]interface{}) {
	zerolog.SetGlobalLevel(getLogLevel(level))

	log.Logger = buildBaseLogger(log.Logger, format).
		With().
		Fields(tags).
		Logger()
}

func Logger(process string, tags map[string]interface{}) zerolog.Logger {
	builder := log.Logger.With()

	if process != "" {
		builder = builder.Str("process", process)
	}

	return builder.Fields(tags).Logger()
}

func LoggerToFile(filename string) (*os.File, error) {
	//nolint:gomnd
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return file, err
	}

	log.Logger = log.Output(file)

	return file, nil
}

func getLogLevel(val string) zerolog.Level {
	level := strings.ToLower(val)

	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

func buildBaseLogger(logger zerolog.Logger, format string) zerolog.Logger {
	logger = logger.With().Str("name", pkg.AppName).Logger()

	switch format {
	case "text":
		//nolint:exhaustruct
		return logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	default:
		return logger
	}
}
