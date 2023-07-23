package support

import (
	"os"
	"strconv"
	"time"
)

// GetEnvString from OS.
func GetEnvString(key, def string) string {
	val := os.Getenv(key)

	if len(val) > 0 {
		return val
	}

	return def
}

// GetEnvInt from OS.
func GetEnvInt(key string, def int) (int, error) {
	val := os.Getenv(key)

	if len(val) > 0 {
		return strconv.Atoi(val)
	}

	return def, nil
}

// GetEnvDur from OS.
func GetEnvDur(key string, def time.Duration) (time.Duration, error) {
	val := os.Getenv(key)

	if len(val) > 0 {
		return time.ParseDuration(val)
	}

	return def, nil
}
