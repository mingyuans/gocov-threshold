package diff

import (
	"bufio"
	"fmt"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/model"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const ConcurrencyLimit = 10 // Limit the number of concurrent goroutines

type Manager struct {
	conf model.Conf
	arg  model.Arg
}

func NewManager(arg model.Arg) *Manager {
	conf, loadConfErr := model.LoadConfFromYAML(arg.ConfPath)
	if loadConfErr != nil {
		log.Get().Fatal("Failed to load configuration", zap.Error(loadConfErr))
	}

	return &Manager{
		conf: conf,
		arg:  arg,
	}
}

func (m Manager) FilterStatements(diffBytes []byte, coveragePath string) ([]model.CovStatement, error) {
	diffBlocks, parseErr := model.LoadPRDiffBlocksFromDiffFile(diffBytes)
	if parseErr != nil {
		return nil, parseErr
	}

	covProfileBytes, readErr := os.ReadFile(coveragePath)
	if readErr != nil {
		return nil, readErr
	}

	covStatements, loadErr := model.LoadCovStatementsFromCovProfile(covProfileBytes)
	if loadErr != nil {
		return nil, loadErr
	}

	filteredStatements, unionErr := m.filterStatementsByPRDiff(diffBlocks, covStatements)
	if unionErr != nil {
		return nil, unionErr
	}
	return m.filterStatementsByPatterns(filteredStatements)
}

func (m Manager) filterStatementsByPRDiff(diffBlocks []model.PRDiffBlock, covStatements []model.CovStatement) ([]model.CovStatement, error) {
	filteredStatements := make([]model.CovStatement, 0)
	for _, diffBlock := range diffBlocks {
		if !m.isFileIncluded(diffBlock.FileName) {
			log.Get().Debug("File excluded from diff", zap.String("fileName", diffBlock.FileName))
			continue
		}

		for _, statement := range covStatements {
			removedModuleFileName := strings.TrimPrefix(statement.FileName, m.arg.Module+"/")
			if removedModuleFileName != diffBlock.FileName {
				continue // Only consider statements from the same file as the diff block
			}

			// Check if statement's range overlaps with diffBlock's range
			if statement.Block.Start <= diffBlock.Block.End && statement.Block.End >= diffBlock.Block.Start {
				filteredStatements = append(filteredStatements, statement)
			}
		}
	}

	return filteredStatements, nil
}

func (m Manager) isFileIncluded(fileName string) bool {
	// Skip _test.go files
	if strings.HasSuffix(fileName, "_test.go") {
		return false
	}

	// Skip non-Go files
	if !strings.HasSuffix(fileName, ".go") {
		return false
	}

	// Skip files in vendor directories
	if strings.Contains(fileName, "vendor/") {
		return true
	}

	// Check if the file is in the exclude list
	for _, excludeDir := range m.conf.Files.Exclude.Dirs {
		if strings.Contains(fileName, excludeDir) {
			log.Get().Debug("File excluded from diff", zap.String("fileName", fileName), zap.String("excludeDir", excludeDir))
			return false
		}
	}

	// Check if the file is matching the include patterns
	for _, excludePattern := range m.conf.Files.Exclude.Patterns {
		fileBase := filepath.Base(fileName)
		matched, err := filepath.Match(excludePattern, fileBase)
		if err != nil {
			log.Get().Error("Invalid exclude pattern", zap.String("pattern", excludePattern), zap.Error(err))
			continue // Skip this pattern if it's invalid
		}
		if matched {
			log.Get().Debug("File excluded from diff", zap.String("fileName", fileName))
			return false
		}
	}

	// Check if the file is in the included directories
	var included = false
	for _, includeDir := range m.conf.Files.Include.Dirs {
		if strings.Contains(fileName, includeDir) {
			included = true
			break
		}
	}

	if len(m.conf.Files.Include.Dirs) > 0 && !included {
		log.Get().Debug("File not included in any include directory", zap.String("fileName", fileName))
		return false
	}

	// Check if the file is matching the include patterns
	included = false
	for _, includePattern := range m.conf.Files.Include.Patterns {
		fileBase := filepath.Base(fileName)
		matched, err := filepath.Match(includePattern, fileBase)
		if err != nil {
			log.Get().Error("Invalid include pattern", zap.String("pattern", includePattern), zap.Error(err))
			continue // Skip this pattern if it's invalid
		}
		if matched {
			included = true
			break
		}
	}

	if len(m.conf.Files.Include.Patterns) > 0 && !included {
		log.Get().Debug("File not included in any include pattern", zap.String("fileName", fileName))
		return false
	}

	return true
}

func (m Manager) filterStatementsByPatterns(covStatements []model.CovStatement) ([]model.CovStatement, error) {
	//Group statements by file name
	groupedStatements := make(map[string][]model.CovStatement)
	for _, statement := range covStatements {
		fileName := strings.TrimPrefix(statement.FileName, m.arg.Module+"/")
		if _, ok := groupedStatements[fileName]; !ok {
			groupedStatements[fileName] = []model.CovStatement{statement}
		}
		groupedStatements[fileName] = append(groupedStatements[fileName], statement)
	}

	// Use a mutex to safely append to the filteredStatements slice across multiple goroutines
	var (
		mu                 sync.Mutex
		filteredStatements []model.CovStatement
		wg                 sync.WaitGroup
		sem                = make(chan struct{}, ConcurrencyLimit) // concurrency limit
	)

	for fileName, statements := range groupedStatements {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore
		go func(fileName string, statements []model.CovStatement) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore

			includedStatements, err := m.filterStatementsForSpecifiedFile(statements, fileName)
			if err != nil {
				log.Get().Error("Failed to filter statements for file", zap.String("fileName", fileName), zap.Error(err))
				return // Skip this file if it cannot be processed
			}

			mu.Lock()
			filteredStatements = append(filteredStatements, includedStatements...)
			mu.Unlock()
		}(fileName, statements)
	}

	wg.Wait()
	return filteredStatements, nil
}

func (m Manager) filterStatementsForSpecifiedFile(statements []model.CovStatement, fileName string) ([]model.CovStatement, error) {
	codeLines, readErr := m.readCodeLinesByCovFileName(fileName)
	if readErr != nil {
		return make([]model.CovStatement, 0), readErr
	}

	if IsFileIgnoreCoverByAnnotation(codeLines) {
		log.Get().Debug("File ignored by annotation",
			zap.String("fileName", fileName))
		return make([]model.CovStatement, 0), nil // Skip this file if it is ignored by annotation
	}

	includedStatements := make([]model.CovStatement, 0)
	for _, statement := range statements {
		// Check if the statement's block is within the code lines
		if statement.Block.Start < 1 || statement.Block.End > len(codeLines) {
			log.Get().Debug("Statement block out of range",
				zap.String("fileName", fileName),
				zap.Int("start", statement.Block.Start),
				zap.Int("end", statement.Block.End))
			continue // Skip this statement if its block is out of range
		}

		// Get the code lines for the statement's block
		statementCodeLines := codeLines[statement.Block.Start-1 : statement.Block.End]
		if IsStatementIgnoreCoverByPatterns(statementCodeLines, m.conf.Statements.Exclude.Patterns) {
			continue
		}
		includedStatements = append(includedStatements, statement)
	}
	return includedStatements, nil
}

func (m Manager) readCodeLinesByCovFileName(fileName string) ([]string, error) {
	trimmedFileName := strings.TrimPrefix(fileName, m.arg.Module+"/")
	codeFilePath := fmt.Sprintf("%s/%s", m.arg.Path, trimmedFileName)

	file, err := os.Open(codeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", codeFilePath, err)
	}
	defer file.Close()

	var codeLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		codeLines = append(codeLines, scanner.Text())
	}
	return codeLines, nil
}
