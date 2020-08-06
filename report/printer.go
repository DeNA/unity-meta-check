package report

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"io"
)

const (
	MissingMetaLinePrefix = "missing "
	DanglingMetaPrefix    = "dangling "
)

func WriteResult(writer io.Writer, result *checker.CheckResult) error {
	if len(result.MissingMeta) > 0 {
		for _, missing := range result.MissingMeta {
			if _, err := io.WriteString(writer, fmt.Sprintf("%s%s\n", MissingMetaLinePrefix, missing)); err != nil {
				return err
			}
		}
	}

	if len(result.DanglingMeta) > 0 {
		for _, dangling := range result.DanglingMeta {
			if _, err := io.WriteString(writer, fmt.Sprintf("%s%s\n", DanglingMetaPrefix, dangling)); err != nil {
				return err
			}
		}
	}

	return nil
}
