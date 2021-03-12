package main

import (
	"gfeed/cmd"
)

var Version string
var Commit string
var BuildDate string

func main() {
	cmd.Execute(cmd.ProcessInfo{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	})
}
