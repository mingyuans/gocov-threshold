package main

import (
	"github.com/mingyuans/gocov-threshold/cmd/threshold/arg"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/pr"
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
