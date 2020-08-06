package options

import (
	"errors"
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
)

type UnityProjectDetector func(unityProj, upmPkg, unityProjSubDir bool, rootDirAbs typedpath.RawPath) (bool, error)

func NewUnityProjectDetector(logger logging.Logger) UnityProjectDetector {
	return func(unityProj, upmPkg, unityProjSubDir bool, rootDirAbs typedpath.RawPath) (bool, error) {
		notUnityProj := upmPkg || unityProjSubDir

		if notUnityProj && unityProj {
			return false, errors.New("must specify one of -upm-package or -unity-project or -unity-project-sub-dir")
		}

		if !notUnityProj && !unityProj {
			logger.Info("none of -upm-package and -unity-project and -unity-project-sub-dir was specified, so try to detect it.")
			isUnityProj, err := hasUnityProjectSpecificDirectory(rootDirAbs, logger)
			if err != nil {
				return false, fmt.Errorf("automatic check mode detection was failed: %q", err.Error())
			}
			return isUnityProj, nil
		}

		return unityProj, nil
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
