package utils

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger of context
type Logger = *logrus.Entry

var logLevel logrus.Level
var logFormater logrus.Formatter

func init() {
	logLevel = getLogLevel()
	logFormater = getLogFormatter()

	logrus.SetLevel(logLevel)
	logrus.SetFormatter(logFormater)
}

// NewLogger create a new logger entry
func NewLogger(context string) Logger {
	return logrus.WithField("context", context)
}

func getLogLevel() logrus.Level {
	level := strings.ToLower(GetEnv("LOG_LEVEL", "info"))

	switch level {
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

func getLogFormatter() logrus.Formatter {
	level := strings.ToLower(GetEnv("LOG_FORMAT", "text"))

	switch level {
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}
