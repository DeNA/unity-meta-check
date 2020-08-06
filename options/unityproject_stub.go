package options

import "github.com/DeNA/unity-meta-check/util/typedpath"

func StubUnityProjectDetector(result bool, err error) UnityProjectDetector {
	return func(_, _, _ bool, _ typedpath.RawPath) (bool, error) {
		return result, err
	}
}
