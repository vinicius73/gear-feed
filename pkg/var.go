package pkg

import (
	"fmt"
	"os"
)

const (
	AppName         = "gfeed"
	maxCommitLength = 8
)

var (
	version   = "dev"
	buildDate = "unknown"
	commit    = "unknown"
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
	if len(commit) > maxCommitLength {
		return commit[:maxCommitLength]
	}

	return commit
}

func Host() string {
	hostname, _ := os.Hostname()

	if hostname == "" {
		return "unknown"
	}

	return hostname
}
