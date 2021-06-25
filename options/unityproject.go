package options

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/pkg/errors"
	"os"
)

type UnityProjectDetector func(rootDirAbs typedpath.RawPath) (checker.TargetType, error)

func NewUnityProjectDetector(logger logging.Logger) UnityProjectDetector {
	return func(rootDirAbs typedpath.RawPath) (checker.TargetType, error) {
		isUnityProj, err := hasUnityProjectSpecificDirectory(rootDirAbs, logger)
		if err != nil {
			return "", errors.Wrap(err, "automatic check mode detection was failed")
		}

		if isUnityProj {
			return checker.TargetTypeIsUnityProjectRootDirectory, nil
		}
		return checker.TargetTypeIsUnityProjectSubDirectory, nil
	}
}

func hasUnityProjectSpecificDirectory(rootDirAbs typedpath.RawPath, logger logging.Logger) (bool, error) {
	assetsDir := rootDirAbs.JoinBaseName(unity.AssetsDirBaseName)
	_, err := os.Stat(string(assetsDir))
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug(fmt.Sprintf("seems a not Unity Project because %s is not found.", assetsDir))
			return false, nil
		}
		return false, err
	}
	logger.Debug(fmt.Sprintf("seems an Unity Project because %s is found.", assetsDir))
	return true, nil
}
