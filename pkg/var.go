package pkg

import (
	"fmt"
	"os"
)

const AppName = "gamer-feed"

var (
	version   string
	buildDate string
	commit    string
)

func Version() string {
	return version
}

func VersionVerbose() string {
	return fmt.Sprintf("Version %s\nRevision %s\nBuild at %s", Version(), Commit(), BuildDate())
}

func BuildDate() string {
	return buildDate
}

func Commit() string {
	return commit
}

func Host() string {
	hostname, _ := os.Hostname()

	if hostname == "" {
		return "unknown"
	}

	return hostname
}
