package report

import (
	"bufio"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io"
	"strings"
)

type Parser func(reader io.Reader) *checker.CheckResult

func NewParser() Parser {
	return func(reader io.Reader) *checker.CheckResult {
		scanner := bufio.NewScanner(reader)

		missingMeta := make([]typedpath.SlashPath, 0)
		danglingMeta := make([]typedpath.SlashPath, 0)

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, MissingMetaLinePrefix) {
				missingMeta = append(missingMeta, typedpath.NewSlashPathUnsafe(strings.TrimPrefix(line, MissingMetaLinePrefix)))
				continue
			}

			if strings.HasPrefix(line, DanglingMetaPrefix) {
				danglingMeta = append(danglingMeta, typedpath.NewSlashPathUnsafe(strings.TrimPrefix(line, DanglingMetaPrefix)))
				continue
			}
		}

		return &checker.CheckResult{
			MissingMeta: missingMeta,
			DanglingMeta: danglingMeta,
		}
	}
}
