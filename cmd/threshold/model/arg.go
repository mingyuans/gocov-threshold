package model

import (
	"fmt"
	"os"
	"strings"
)

type Arg struct {
	Module              string
	Threshold           float64
	Path                string
	Coverprofile        string
	LoggerLevel         string
	GithubToken         string
	ConfPath            string
	PrintUncoveredLines bool
}

// ParseArg parses command-line flags into an Arg struct.
func ParseArg() Arg {
	var a Arg
	a.Module = getActionInput("module")
	a.Threshold = 0
	if threshold := getActionInput("threshold"); threshold != "" {
		_, scanErr := fmt.Sscanf(threshold, "%f", &a.Threshold)
		if scanErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error parsing threshold: %v\n", scanErr)
			os.Exit(1)
		}
	}
	a.Path = getActionInput("path")
	a.Coverprofile = getActionInput("coverprofile")
	a.LoggerLevel = getActionInput("logger-level")
	a.GithubToken = getActionInput("token")
	a.ConfPath = getActionInput("conf")
	printUncoveredLines := getActionInput("print-uncovered-lines")
	a.PrintUncoveredLines = strings.ToLower(printUncoveredLines) == "true"
	return a
}

func getActionInput(input string) string {
	return os.Getenv(
		fmt.Sprintf("INPUT_%s", strings.ToUpper(input)),
	)
}
