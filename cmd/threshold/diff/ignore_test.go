package diff

import (
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsFileIgnoreCover(t *testing.T) {
	t.Run("ignore coverage", func(t *testing.T) {
		fileLines := []string{
			"// gocover:ignore",
			"package main",
			"",
			"func CalculateSum(a, b int) int {",
			"	return a + b",
			"}",
		}
		ignore := IsFileIgnoreCoverByAnnotation(fileLines)
		require.True(t, ignore, "Expected file to be ignored for coverage")
	})

	t.Run("wrong annotation", func(t *testing.T) {
		fileLines := []string{
			"// gocover:ignoretest",
			"package main",
			"",
			"func CalculateSum(a, b int) int {",
			"	return a + b",
			"}",
		}
		ignore := IsFileIgnoreCoverByAnnotation(fileLines)
		require.False(t, ignore, "Expected file not to be ignored for coverage")
	})
}

func TestIsStatementIgnoreCover(t *testing.T) {
	log.Init("debug")

	t.Run("ignore coverage", func(t *testing.T) {
		fileLines := []string{
			"func CalculateSum(a, b int) int {",
			"	// gocover:ignore",
			"	return a + b",
			"}",
		}
		ignore := IsStatementIgnoreCoverByPatterns(fileLines, nil)
		require.True(t, ignore, "Expected file to be ignored for coverage")
	})

	t.Run("wrong annotation", func(t *testing.T) {
		fileLines := []string{
			"func CalculateSum(a, b int) int {",
			"	// gocover:ignoressssss",
			"	return a + b",
			"}",
		}
		ignore := IsStatementIgnoreCoverByPatterns(fileLines, nil)
		require.False(t, ignore, "Expected file not to be ignored for coverage")
	})
}
