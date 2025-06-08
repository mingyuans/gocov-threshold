package diff

import (
	"fmt"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/model"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestManager_union(t *testing.T) {
	log.Init("debug")

	diffBytes := loadTestData(t, "example-diff.diff")
	covProfileBytes := loadTestData(t, "example-overage.out")
	diffBlocks, diffErr := model.LoadPRDiffBlocksFromDiffFile(diffBytes)
	require.NoError(t, diffErr)
	require.NotEmpty(t, diffBlocks)
	covStatements, covErr := model.LoadCovStatementsFromCovProfile(covProfileBytes)
	require.NoError(t, covErr)
	require.NotEmpty(t, covStatements)

	t.Run("filterStatementsByPRDiff", func(t *testing.T) {
		m := Manager{
			arg:  model.Arg{Module: "mingyuans/gocov-threshold"},
			conf: model.Conf{},
		}

		_, err := m.filterStatementsByPRDiff(diffBlocks, covStatements)
		require.NoError(t, err)
	})

	t.Run("union2", func(t *testing.T) {
		m := Manager{
			arg:  model.Arg{Module: "mingyuans/gocov-threshold"},
			conf: model.Conf{},
		}

		tempDiffBlocks := make([]model.PRDiffBlock, 0)
		tempDiffBlocks = append(tempDiffBlocks, model.PRDiffBlock{
			FileName: "example/main.go",
			Block: model.Block{
				Start: 10,
				End:   15,
			},
		})

		tempCovStatements := make([]model.CovStatement, 0)
		tempCovStatements = append(tempCovStatements, model.CovStatement{
			FileName: "mingyuans/gocov-threshold/example/main.go",
			Block:    model.Block{Start: 10, End: 15},
		})
		tempCovStatements = append(tempCovStatements, model.CovStatement{
			FileName: "mingyuans/gocov-threshold/example/cmd.go",
			Block:    model.Block{Start: 11, End: 12},
		})

		statements, err := m.filterStatementsByPRDiff(tempDiffBlocks, tempCovStatements)
		require.NoError(t, err)
		require.Len(t, statements, 1, "Expected only one statement to match the diff block")
		require.Equal(t, statements[0].FileName, "mingyuans/gocov-threshold/example/main.go")
	})

}

func loadTestData(t *testing.T, fileName string) []byte {
	t.Helper()
	path := "../../../testdata/" + fileName
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data from %s: %v", fileName, err)
	}
	return data
}

func TestManager_isFileIncluded(t *testing.T) {
	log.Init("debug")

	t.Run("nonGoFiles", func(t *testing.T) {
		m := Manager{}
		fileName := "testdata/example.txt"
		isIncluded := m.isFileIncluded(fileName)
		require.False(t, isIncluded, "Expected file to be excluded")
	})

	t.Run("includedDirs", func(t *testing.T) {
		m := Manager{conf: model.Conf{}}
		// Only include Go files under the "example" directory
		m.conf.Files.Include.Dirs = []string{"example"}
		fileName := "testdata/main.go"
		isIncluded := m.isFileIncluded(fileName)
		require.False(t, isIncluded, "Expected file to be excluded")
	})

	t.Run("excludeDirs", func(t *testing.T) {
		m := Manager{conf: model.Conf{}}
		// Only include Go files under the "example" directory
		m.conf.Files.Exclude.Dirs = []string{"testdata"}
		fileName := "testdata/main.go"
		isIncluded := m.isFileIncluded(fileName)
		require.False(t, isIncluded, "Expected file to be excluded")
	})

	t.Run("includedPatterns", func(t *testing.T) {
		m := Manager{conf: model.Conf{}}
		m.conf.Files.Include.Patterns = []string{"*_model.go"}
		fileName := "testdata/main.go"
		isIncluded := m.isFileIncluded(fileName)
		require.False(t, isIncluded, "Expected file to be excluded")

		fileName = "testdata/main_model.go"
		isIncluded = m.isFileIncluded(fileName)
		require.True(t, isIncluded, "Expected file to be included")
	})

	t.Run("excludePatterns", func(t *testing.T) {
		m := Manager{conf: model.Conf{}}
		// Only include Go files under the "example" directory
		m.conf.Files.Exclude.Patterns = []string{"*_model.go"}
		fileName := "testdata/main.go"
		isIncluded := m.isFileIncluded(fileName)
		require.True(t, isIncluded, "Expected file to be included")

		fileName = "testdata/main_model.go"
		isIncluded = m.isFileIncluded(fileName)
		require.False(t, isIncluded, "Expected file to be excluded")

	})
}

func TestManager_filterStatementsForSpecifiedFile(t *testing.T) {
	log.Init("debug")

	t.Run("ignore file", func(t *testing.T) {
		pwd := os.Getenv("PWD")
		m := Manager{
			arg: model.Arg{
				Module: "mingyuans/gocov-threshold",
				Path:   fmt.Sprintf("%s/../../..", pwd),
			},
		}

		statements := []model.CovStatement{}
		statements = append(statements,
			model.CovStatement{
				FileName: "mingyuans/gocov-threshold/example/ignore.go",
				Block: model.Block{
					Start: 10,
					End:   15,
				},
			})

		filteredStatements, err := m.filterStatementsForSpecifiedFile(statements, "mingyuans/gocov-threshold/example/ignore.go")
		require.NoError(t, err)
		require.Empty(t, filteredStatements, "Expected no statements to be returned for ignored file")
	})

	t.Run("ignore statement by gocover:ignore", func(t *testing.T) {
		pwd := os.Getenv("PWD")
		m := Manager{
			arg: model.Arg{
				Module: "github.com/mingyuans/gocov-threshold",
				Path:   fmt.Sprintf("%s/../../..", pwd),
			},
		}

		statements := []model.CovStatement{}
		statements = append(statements,
			model.CovStatement{
				FileName: "github.com/mingyuans/gocov-threshold/example/main.go",
				Block: model.Block{
					Start: 14,
					End:   16,
				},
			},
			model.CovStatement{
				FileName: "github.com/mingyuans/gocov-threshold/example/main.go",
				Block: model.Block{
					Start: 27,
					End:   30,
				},
			})

		filteredStatements, err := m.filterStatementsForSpecifiedFile(statements, "github.com/mingyuans/gocov-threshold/example/main.go")
		require.NoError(t, err)
		require.Len(t, filteredStatements, 1, "Expected only one statement to match the diff block")
		require.Equal(t, 14, filteredStatements[0].Block.Start, "Expected the start of the block to be 14")
	})

	t.Run("ignore statement by global custom regex", func(t *testing.T) {
		pwd := os.Getenv("PWD")
		m := Manager{
			arg: model.Arg{
				Module: "mingyuans/gocov-threshold",
				Path:   fmt.Sprintf("%s/../../..", pwd),
			},
			conf: model.Conf{
				Statements: model.StatementConf{
					Exclude: model.StatementPattern{
						Patterns: []string{".*mu.Lock.*"},
					},
				},
			},
		}

		statements := []model.CovStatement{}
		statements = append(statements,
			model.CovStatement{
				FileName: "mingyuans/gocov-threshold/example/main.go",
				Block: model.Block{
					Start: 14,
					End:   16,
				},
			},
			model.CovStatement{
				FileName: "mingyuans/gocov-threshold/example/main.go",
				Block: model.Block{
					Start: 33,
					End:   40,
				},
			})

		filteredStatements, err := m.filterStatementsForSpecifiedFile(statements, "mingyuans/gocov-threshold/example/main.go")
		require.NoError(t, err)
		require.Len(t, filteredStatements, 1, "Expected only one statement to match the diff block")
		require.Equal(t, 14, filteredStatements[0].Block.Start, "Expected the start of the block to be 14")
	})
}

func TestManager_FilterStatements(t *testing.T) {
	log.Init("debug")

	t.Run("basic filter", func(t *testing.T) {
		pwd := os.Getenv("PWD")
		m := Manager{
			arg: model.Arg{
				Module: "github.com/mingyuans/gocov-threshold",
				Path:   fmt.Sprintf("%s/../../..", pwd),
			},
			conf: model.Conf{
				Files: model.FileConf{
					Include: model.FilePattern{
						Dirs: []string{"example"},
					},
				},
			},
		}

		diffBytes := loadTestData(t, "example-diff.diff")
		coveragePath := fmt.Sprintf("%s/../../../testdata/example-overage.out", pwd)
		statements, err := m.FilterStatements(diffBytes, coveragePath)
		require.NoError(t, err)
		require.NotEmpty(t, statements, "Expected filtered statements to be non-empty")
	})

	t.Run("empty diff", func(t *testing.T) {
		pwd := os.Getenv("PWD")
		m := Manager{
			arg: model.Arg{
				Module:   "mingyuans/gocov-threshold",
				Path:     fmt.Sprintf("%s/../../..", pwd),
				ConfPath: "../../../testdata/test-conf.yaml",
			},
			conf: model.Conf{},
		}

		diffBytes := []byte{}
		coveragePath := "../../../testdata/example-overage.out"

		statements, err := m.FilterStatements(diffBytes, coveragePath)
		require.NoError(t, err)
		require.Empty(t, statements)
	})
}
