package checker

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/pathutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"sort"
)

type CheckingWorker func(ignoreCase bool, reader <-chan typedpath.SlashPath) (*CheckResult, error)

func NewCheckingWorker(requiresMeta unity.MetaNecessity, logger logging.Logger) CheckingWorker {
	return func(ignoreCase bool, reader <-chan typedpath.SlashPath) (*CheckResult, error) {
		expectedMetaSet := pathutil.NewPathSet(ignoreCase)
		// NOTE: This set hold all exist files as elements (without .meta).
		//       So, the elements may not effective .meta, but it should be allowed.
		// WHY: See chkworker_test.go
		allowedMetaSet := pathutil.NewPathSet(ignoreCase)
		actualMetaSet := pathutil.NewPathSet(ignoreCase)

		for targetPath := range reader {
			if unity.IsMeta(targetPath) {
				logger.Debug(fmt.Sprintf("meta found ... %s", targetPath))
				actualMetaSet.Add(targetPath)
			} else if requiresMeta(targetPath) {
				logger.Debug(fmt.Sprintf("needs meta ... %s", targetPath))
				expectedMetaSet.Add(unity.MetaPath(targetPath))
				allowedMetaSet.Add(unity.MetaPath(targetPath))
			} else {
				allowedMetaSet.Add(unity.MetaPath(targetPath))
				logger.Debug(fmt.Sprintf("skipped ... %s", targetPath))
			}
		}

		missingMeta := expectedMetaSet.Difference(actualMetaSet)
		sort.Slice(missingMeta, func(i, j int) bool {
			return missingMeta[i] < missingMeta[j]
		})

		danglingMeta := actualMetaSet.Difference(allowedMetaSet)
		sort.Slice(danglingMeta, func(i, j int) bool {
			return danglingMeta[i] < danglingMeta[j]
		})

		return NewCheckResult(missingMeta, danglingMeta), nil
	}
}
