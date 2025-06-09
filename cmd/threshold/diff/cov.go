package diff

import (
	"bufio"
	"fmt"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/model"
	"os"
	"strings"
)

func (m Manager) getCovStatement(statement model.CovStatement) ([]string, error) {
	// Remove the module prefix from the file name
	fileName := strings.TrimPrefix(statement.FileName, m.arg.Module+"/")
	codeFilePath := fmt.Sprintf("%s/%s", m.arg.Path, fileName)

	file, err := os.Open(codeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", codeFilePath, err)
	}
	defer file.Close()

	var codeLines []string
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		if lineNum >= statement.Hunk.Start && lineNum <= statement.Hunk.End {
			codeLines = append(codeLines, scanner.Text())
		}
		if lineNum > statement.Hunk.End {
			break
		}
		lineNum++
	}
	//statement.Code = strings.Join(codeLines, "\n")
	return codeLines, nil
}
