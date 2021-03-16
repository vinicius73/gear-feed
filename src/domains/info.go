package domains

var version string
var commit string
var buildDate string

var info ProcessInfo

type ProcessInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

func init() {
	info = ProcessInfo{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}
}

func Info() ProcessInfo {
	return info
}
