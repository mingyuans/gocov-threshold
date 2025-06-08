package diff

import (
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"go.uber.org/zap"
	"regexp"
)

var IgnoreRegex = regexp.MustCompile(`.*//\s*(gocover:ignore)\s*$`)

func IsFileIgnoreCoverByAnnotation(fileLines []string) bool {
	if len(fileLines) == 0 {
		return true
	}
	match := IgnoreRegex.FindString(fileLines[0])
	return match != ""
}

func IsStatementIgnoreCoverByPatterns(lines []string, patterns []string) bool {
	regexes := make([]*regexp.Regexp, 0)
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			log.Get().Warn("Failed to compile regex pattern", zap.String("pattern", pattern), zap.Error(err))
			continue // skip invalid patterns
		}
		regexes = append(regexes, re)
	}

	regexes = append(regexes, IgnoreRegex)
	for _, line := range lines {
		for _, re := range regexes {
			if re.MatchString(line) {
				log.Get().Debug("ignore cov by pattern",
					zap.String("line", line),
					zap.String("pattern", re.String()))
				return true // Skip this statement if it matches any exclude pattern
			}
		}
	}

	return false
}
