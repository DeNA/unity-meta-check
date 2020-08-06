package checker

import "github.com/DeNA/unity-meta-check/util/typedpath"

type CheckResult struct {
	MissingMeta  []typedpath.SlashPath
	DanglingMeta []typedpath.SlashPath
}

func (c CheckResult) Empty() bool {
	return c.Len() == 0
}

func (c CheckResult) Len() int {
	return len(c.MissingMeta) + len(c.DanglingMeta)
}

func NewCheckResult(missingMeta []typedpath.SlashPath, danglingMeta []typedpath.SlashPath) *CheckResult {
	return &CheckResult{
		MissingMeta:  missingMeta,
		DanglingMeta: danglingMeta,
	}
}
