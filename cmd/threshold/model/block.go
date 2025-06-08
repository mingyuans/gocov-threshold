package model

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Block struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type CovStatement struct {
	FileName       string `json:"file_name"`
	Block          Block  `json:"block"`
	StatementCount int    `json:"statement_count"`
	ExecutionCount int    `json:"execution_count"`
	// Only filled after filtering by patterns
	CodeLines []string `json:"code_lines"`
}

type PRDiffBlock struct {
	FileName string `json:"file_name"`
	Block    Block  `json:"block"`
}

func LoadCovStatementsFromCovProfile(covFileBytes []byte) ([]CovStatement, error) {
	var covStatements []CovStatement
	scanner := bufio.NewScanner(bytes.NewReader(covFileBytes))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue // skip header
		}
		// Example line: path/to/file.go:10.12,12.3 2 1
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		fileAndBlock := parts[0]
		statementCountStr := parts[1]
		executionCountStr := parts[2]

		fileBlockParts := strings.Split(fileAndBlock, ":")
		if len(fileBlockParts) != 2 {
			continue
		}
		fileName := fileBlockParts[0]
		blockRange := fileBlockParts[1]

		blockParts := strings.Split(blockRange, ",")
		if len(blockParts) != 2 {
			continue
		}
		startStr := blockParts[0]
		endStr := blockParts[1]

		startLineCol := strings.Split(startStr, ".")
		endLineCol := strings.Split(endStr, ".")
		if len(startLineCol) != 2 || len(endLineCol) != 2 {
			continue
		}
		startLine, err1 := strconv.Atoi(startLineCol[0])
		endLine, err2 := strconv.Atoi(endLineCol[0])
		if err1 != nil || err2 != nil {
			continue
		}

		statementCount, err3 := strconv.Atoi(statementCountStr)
		executionCount, err4 := strconv.Atoi(executionCountStr)
		if err3 != nil || err4 != nil {
			continue
		}

		covStatements = append(covStatements, CovStatement{
			FileName: fileName,
			Block: Block{
				Start: startLine,
				End:   endLine,
			},
			StatementCount: statementCount,
			ExecutionCount: executionCount,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan coverage profile: %w", err)
	}
	return covStatements, nil
}

func LoadPRDiffBlocksFromDiffFile(diffFileBytes []byte) ([]PRDiffBlock, error) {
	var diffBlocks []PRDiffBlock
	scanner := bufio.NewScanner(bytes.NewReader(diffFileBytes))
	var currentFile string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "diff --git") {
			// New file section starts
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				currentFile = strings.TrimPrefix(parts[2], "a/")
			}
			continue
		}
		if strings.HasPrefix(line, "@@") {
			// Block of code starts
			blockPartsInNewFile := strings.Split(line, " ")
			if len(blockPartsInNewFile) < 2 {
				continue
			}
			blockRangeInNewFile := blockPartsInNewFile[2]
			blockPartsInNewFile = strings.Split(blockRangeInNewFile, ",")
			if len(blockPartsInNewFile) < 2 {
				continue
			}
			startLineStr := strings.TrimPrefix(blockPartsInNewFile[0], "+")
			startLine, err := strconv.Atoi(startLineStr)
			if err != nil {
				continue
			}
			endLineStr := blockPartsInNewFile[1]
			endLine, err := strconv.Atoi(endLineStr)
			if err != nil {
				continue
			}
			diffBlocks = append(diffBlocks, PRDiffBlock{
				FileName: currentFile,
				Block: Block{
					Start: startLine,
					End:   startLine + endLine - 1,
				},
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan diff file: %w", err)
	}
	return diffBlocks, nil
}
