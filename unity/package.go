package unity

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const PackagesDirname typedpath.BaseName = "Packages"
const LocalPkgPrefix = "file:"
const LocalPkgPrefixLen = len(LocalPkgPrefix)

type FindPackages func(projRoot typedpath.RawPath) ([]*FoundPackage, error)

// NewFindPackages returns the dependencies field of manifest.json. For example,
// 	Input:
// 	"dependencies": {
// 	  ...
// 	  "com.my.pkg": "1.0.0"
// 	  "com.my.local.pkg": "file:../MyLocalPkg/com.my.local.pkg"
// 	  "com.my.another.local.pkg": "file:../MyLocalPkg/com.my.another.local.pkg"
// 	  ...
// 	}
//
// 	Output: []string{"Packages/com.my.pkg", "MyLocalPkg/com.my.local.pkg", "MyLocalPkg/com.my.another.local.pkg"}
func NewFindPackages(logger logging.Logger) FindPackages {
	return func(rootDirAbs typedpath.RawPath) ([]*FoundPackage, error) {
		if !rootDirAbs.IsAbs() {
			return nil, fmt.Errorf("project root path must be absolute: %q", rootDirAbs)
		}

		packagesDirAbsPath := rootDirAbs.JoinBaseName(PackagesDirname)
		manifestAbsPath := packagesDirAbsPath.JoinBaseName(ManifestBasename)

		logger.Debug(fmt.Sprintf("Reading manifest.json at: %q", manifestAbsPath))
		manifest, err := ReadManifest(manifestAbsPath)
		if err != nil {
			return nil, err
		}

		result := make([]*FoundPackage, 0)
		for key, value := range manifest.Dependencies {
			if !strings.HasPrefix(value, LocalPkgPrefix) {
				relPath := typedpath.NewRawPath(PackagesDirname, typedpath.BaseName(key))
				absPath := rootDirAbs.JoinRawPath(relPath)
				if _, err := os.Stat(string(absPath)); err != nil {
					if os.IsNotExist(err) {
						logger.Info(fmt.Sprintf("skip package %q because it does not exist: %s", key, relPath))
						continue
					}
					return nil, err
				}

				result = append(result, &FoundPackage{
					FilePrefix: false,
					AbsPath:    absPath,
					RelPath:    relPath,
				})
				logger.Info(fmt.Sprintf("package %q found: %s", key, relPath))
				continue
			}

			// NOTE: This relative path is from Packages/manifest.json.
			localPkgAbsRawPath := packagesDirAbsPath.JoinRawPath(typedpath.SlashPath(value[LocalPkgPrefixLen:]).ToRaw())
			localPkgRelRawPath, err := rootDirAbs.Rel(localPkgAbsRawPath)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot detect Local package location")
			}
			if _, err := os.Stat(string(localPkgAbsRawPath)); err != nil {
				if os.IsNotExist(err) {
					logger.Info(fmt.Sprintf("skip local package %q because it does not exist: %s", key, localPkgRelRawPath))
					continue
				}
				return nil, err
			}

			result = append(result, &FoundPackage{
				FilePrefix: true,
				AbsPath:    localPkgAbsRawPath,
				RelPath:    localPkgRelRawPath,
			})
			logger.Info(fmt.Sprintf("local package %q found: %s", key, localPkgRelRawPath))
		}
		return result, nil
	}
}

type FoundPackage struct {
	FilePrefix bool
	AbsPath    typedpath.RawPath
	RelPath    typedpath.RawPath
}

func FoundPackagesToSlashRelPaths(foundPackages []*FoundPackage) []typedpath.SlashPath {
	result := make([]typedpath.SlashPath, len(foundPackages))
	for i, foundPkg := range foundPackages {
		result[i] = foundPkg.RelPath.ToSlash()
	}
	return result
}
