package main

import (
	"fmt"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/diff"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/model"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/pr"
	"go.uber.org/zap"
)

const DiffFilePath = "pr.diff"

func main() {
	actionArg := model.ParseArg()
	log.Init(actionArg.LoggerLevel)
	log.Get().Debug("Arguments parsed", zap.Any("arg", actionArg))

	diffManager := diff.NewManager(actionArg)

	prService := pr.NewService(actionArg)
	log.Get().Debug("PR Service initialized", zap.Any("env", prService.GetEnvironment()),
		zap.Any("pr", prService.GetPRInfo()))

	diffBytes, downloadErr := prService.DownloadDiff()
	if downloadErr != nil {
		panic(fmt.Sprintf("Failed to download and save PR diff: %s", downloadErr.Error()))
	}

	diffManager.FilterStatements(diffBytes, actionArg.CoveragePath)

}
