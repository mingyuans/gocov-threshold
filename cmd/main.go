package main

import (
	"github.com/mingyuans/gocov-threshold/cmd/arg"
	"github.com/mingyuans/gocov-threshold/cmd/log"
	"github.com/mingyuans/gocov-threshold/cmd/pr"
	"go.uber.org/zap"
)

func main() {
	actionArg := arg.ParseArg()
	log.Init(actionArg.LoggerLevel)
	log.Get().Debug("Arguments parsed", zap.Any("arg", actionArg))

	prService := pr.NewService()
	log.Get().Debug("PR Service initialized", zap.Any("env", prService.GetEnvironment()),
		zap.Any("pr", prService.GetPRInfo()))
}
