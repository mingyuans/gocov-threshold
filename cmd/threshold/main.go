package main

import (
	"fmt"
	"github.com/actions-go/toolkit/core"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/diff"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/model"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/pr"
	"go.uber.org/zap"
	"strings"
)

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
		log.Get().Debug("Statements to calculate", zap.Any("statement", statement))
		if statement.ExecutionCount > 0 {
			coveredStatements++
		} else if actionArg.PrintUncoveredLines {
			core.Infof("Uncovered lines:\n%s", strings.Join(statement.CodeLines, "\n"))
		}
	}
	var coverage = 100.0
	if totalStatements != 0 {
		coverage = float64(coveredStatements) * 100 / float64(totalStatements)
	}
	log.Get().Info("Coverage statistics",
		zap.Int("total_statements", totalStatements),
		zap.Int("covered_statements", coveredStatements),
		zap.Float64("coverage_percentage", coverage))

	if coverage < actionArg.Threshold {
		core.SetFailedf("Coverage %.2f%% is below the threshold %.2f%%", coverage, actionArg.Threshold)
	} else {
		log.Get().Info("Coverage meets the threshold",
			zap.Float64("coverage", coverage),
			zap.Float64("threshold", actionArg.Threshold))
	}
	core.SetOutput("gocov", fmt.Sprintf("%.2f", coverage))
}
