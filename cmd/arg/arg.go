package arg

import "flag"

type Arg struct {
	IgnoreMain   bool
	Module       string
	Threshold    int
	Path         string
	Coverprofile string
	LoggerLevel  string
}

// ParseArg parses command-line flags into an Arg struct.
func ParseArg() Arg {
	var a Arg
	flag.BoolVar(&a.IgnoreMain, "ignore-main", false, "Ignore main package")
	flag.StringVar(&a.Module, "module", "", "Module name")
	flag.IntVar(&a.Threshold, "threshold", 0, "Threshold value")
	flag.StringVar(&a.Path, "path", "", "Path to target")
	flag.StringVar(&a.Coverprofile, "coverprofile", "", "Coverprofile output path")
	flag.StringVar(&a.LoggerLevel, "logger-level", "", "Logger level (debug, info, warn, error)")
	flag.Parse()
	return a
}
