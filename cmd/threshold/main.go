package main

import (
	"fmt"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/arg"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/pr"
	"go.uber.org/zap"
)

func main() {
	actionArg := arg.ParseArg()
	log.Init(actionArg.LoggerLevel)
	log.Get().Debug("Arguments parsed", zap.Any("arg", actionArg))
	fmt.Printf("Parsed arguments: %+v\n", actionArg)

	prService := pr.NewService(actionArg)
	log.Get().Debug("PR Service initialized", zap.Any("env", prService.GetEnvironment()),
		zap.Any("pr", prService.GetPRInfo()))

	downloadErr := prService.DownloadAndSaveDiff("pr.diff")
	if downloadErr != nil {
		panic(fmt.Sprintf("Failed to download and save PR diff: %s", downloadErr.Error()))
	}
}
