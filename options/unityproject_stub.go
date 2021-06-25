package options

import (
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubUnityProjectDetector(result checker.TargetType, err error) UnityProjectDetector {
	return func(_ typedpath.RawPath) (checker.TargetType, error) {
		return result, err
	}
}
