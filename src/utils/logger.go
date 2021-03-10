package utils

import (
	"github.com/sirupsen/logrus"
)

// Logger of context
type Logger = *logrus.Entry

func init() {
	// logrus.SetLevel(logrus.DebugLevel)
}

// NewLogger create a new logger entry
func NewLogger(context string) Logger {
	l := logrus.New()
	// l.SetLevel(logrus.DebugLevel)

	return l.WithField("context", context)
}
