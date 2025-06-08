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

	statements, err := diffManager.FilterStatements(diffBytes, actionArg.Coverprofile)
	if err != nil {
		log.Get().Fatal("Failed to filter statements", zap.Error(err))
	}

	var totalStatements = len(statements)
	var coveredStatements = 0
	for _, statement := range statements {
		log.Get().Debug("Statement", zap.Any("statement", statement))
		if statement.ExecutionCount > 0 {
			coveredStatements++
		}
	}
	coverage := float64(coveredStatements) * 100 / float64(totalStatements)
	log.Get().Info("Coverage statistics",
		zap.Int("total_statements", totalStatements),
		zap.Int("covered_statements", coveredStatements),
		zap.Float64("coverage_percentage", coverage))
}
