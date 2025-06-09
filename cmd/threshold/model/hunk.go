package model

import (
	"bufio"
	"bytes"
	"fmt"
	godiff "github.com/sourcegraph/go-diff/diff"
	"strconv"
	"strings"
)

type Hunk struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type CovStatement struct {
	FileName       string `json:"file_name"`
	Hunk           Hunk   `json:"hunk"`
	StatementCount int    `json:"statement_count"`
	ExecutionCount int    `json:"execution_count"`
	// Only filled after filtering by patterns
	CodeLines []string `json:"code_lines"`
}

type PRHunk struct {
	FileName string `json:"file_name"`
	Hunk     Hunk   `json:"hunk"`
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
			Hunk: Hunk{
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

func LoadPRDiffBlocksFromDiffFile(diffFileBytes []byte) ([]PRHunk, error) {
	var prHunks []PRHunk
	fileDiff, parseErr := godiff.ParseMultiFileDiff(diffFileBytes)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse diff fd: %w", parseErr)
	}
	for _, fd := range fileDiff {
		newFileName := strings.TrimPrefix(fd.NewName, "b/")
		for _, hunk := range fd.Hunks {
			var unifiedOffset = getUnifiedOffsetFromHunk(*hunk)
			blockStart := hunk.NewStartLine + unifiedOffset
			prHunk := PRHunk{
				FileName: newFileName,
				Hunk: Hunk{
					Start: int(blockStart),
					End:   int(blockStart + hunk.NewLines - unifiedOffset),
				},
			}
			prHunks = append(prHunks, prHunk)
		}
	}
	return prHunks, nil
}

func getUnifiedOffsetFromHunk(hunk godiff.Hunk) int32 {
	var unifiedOffset int32 = 0
	lines := bytes.Split(hunk.Body, []byte{'\n'})
	for _, line := range lines {
		if len(line) == 0 || (line[0] != '+' && line[0] != '-') {
			unifiedOffset++
			continue
		}
		// If we encounter a line that starts with '+' or '-', we stop counting
		break
	}
	return unifiedOffset
}
