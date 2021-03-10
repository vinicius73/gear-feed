package utils

import (
	"github.com/sirupsen/logrus"
)

// Logger of context
type Logger = *logrus.Entry

// NewLogger create a new logger entry
func NewLogger(context string) Logger {
	return logrus.New().WithField("context", context)
}
