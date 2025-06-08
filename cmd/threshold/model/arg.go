package model

import (
	"fmt"
	"os"
	"strings"
)

type Arg struct {
	IgnoreMain   bool
	Module       string
	Threshold    int
	Path         string
	Coverprofile string
	LoggerLevel  string
	GithubToken  string
	ConfPath     string
}

// ParseArg parses command-line flags into an Arg struct.
func ParseArg() Arg {
	var a Arg
	a.IgnoreMain = getActionInput("ignore-main") == "true"
	a.Module = getActionInput("module")
	a.Threshold = 80
	if threshold := getActionInput("threshold"); threshold != "" {
		_, _ = fmt.Sscanf(threshold, "%d", &a.Threshold)
	}
	a.Path = getActionInput("path")
	a.Coverprofile = getActionInput("coverprofile")
	a.LoggerLevel = getActionInput("logger-level")
	a.GithubToken = getActionInput("token")
	a.ConfPath = getActionInput("conf")
	return a
}

func getActionInput(input string) string {
	return os.Getenv(
		fmt.Sprintf("INPUT_%s", strings.ToUpper(input)),
	)
}
