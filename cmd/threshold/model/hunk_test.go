package model

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadCovStatementsFromCovProfile(t *testing.T) {
	pwdDir, wdErr := os.Getwd()
	require.NoError(t, wdErr)
	covProfile := fmt.Sprintf("%s/../../../testdata/example-overage.out", pwdDir)

	t.Run("LoadCovStatementsFromCovProfile", func(t *testing.T) {
		covProfileBytes, readErr := os.ReadFile(covProfile)
		require.NoError(t, readErr)
		require.NotEmpty(t, covProfileBytes)
		statements, loadErr := LoadCovStatementsFromCovProfile(covProfileBytes)
		require.NoError(t, loadErr)
		require.NotEmpty(t, statements)
		require.NotEmpty(t, statements[0].FileName)
	})
}

func TestLoadPRDiffBlocksFromDiffFile(t *testing.T) {
	pwdDir, wdErr := os.Getwd()
	require.NoError(t, wdErr)

	t.Run("LoadPRDiffBlocksFromDiffFile", func(t *testing.T) {
		diffFile := fmt.Sprintf("%s/../../../testdata/example-diff.diff", pwdDir)
		diffBytes, readErr := os.ReadFile(diffFile)
		require.NoError(t, readErr)
		require.NotEmpty(t, diffBytes)
		blocks, loadErr := LoadPRDiffBlocksFromDiffFile(diffBytes)
		require.NoError(t, loadErr)
		require.NotEmpty(t, blocks)
		require.NotEmpty(t, blocks[0].FileName)
	})

	t.Run("LoadPRDiffBlocksFromDiffFile", func(t *testing.T) {
		diffFile := fmt.Sprintf("%s/../../../testdata/example-2.diff", pwdDir)
		diffBytes, readErr := os.ReadFile(diffFile)
		require.NoError(t, readErr)
		require.NotEmpty(t, diffBytes)
		blocks, loadErr := LoadPRDiffBlocksFromDiffFile(diffBytes)
		require.NoError(t, loadErr)
		require.NotEmpty(t, blocks)
		require.NotEmpty(t, blocks[0].FileName)
	})
}

//func TestLoadPRBlocksFromDiffFile(t *testing.T) {
//	log.Init("debug")
//	pwdDir, wdErr := os.Getwd()
//	require.NoError(t, wdErr)
//	diffFile := fmt.Sprintf("%s/../../../testdata/example-2.diff", pwdDir)
//
//	t.Run("LoadPRDiffBlocksFromDiffFile", func(t *testing.T) {
//		diffBytes, readErr := os.ReadFile(diffFile)
//		require.NoError(t, readErr)
//		require.NotEmpty(t, diffBytes)
//		blocks, loadErr := LoadPRBlocksFromDiffFile(diffBytes)
//		require.NoError(t, loadErr)
//		require.NotEmpty(t, blocks)
//		require.NotEmpty(t, blocks[0].FileName)
//	})
//}
