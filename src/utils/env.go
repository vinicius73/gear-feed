package utils

import "os"

// GetEnv from OS
func GetEnv(key, def string) string {
	val := os.Getenv(key)

	if len(val) > 0 {
		return val
	}

	return def
}
