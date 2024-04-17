package unity

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/pathutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"strings"
)

type MetaNecessity func(targetPath typedpath.SlashPath) bool

const AssetsDirBaseName typedpath.BaseName = "Assets"

func NewMetaNecessityInUnityProject(pkgPaths []typedpath.SlashPath) MetaNecessity {
	localPkgTree := pathutil.NewPathTree(pkgPaths...)
	requiresMetaInUpmPkg := NewMetaNecessityInUnityProjectSubDir()
	manifestPath := typedpath.SlashPathFromBaseName(PackagesDirname).JoinBaseName(ManifestBasename)

	return func(targetPath typedpath.SlashPath) bool {
		elements := pathutil.SplitPathElements(targetPath)
		if len(elements) == 0 {
			return false
		}

		firstElement := elements[0]
		if firstElement != AssetsDirBaseName {
			if localPkgTree.Member(elements) {
				return requiresMetaInUpmPkg(targetPath)
			}
			return false
		}

		if targetPath == manifestPath || (firstElement == AssetsDirBaseName && len(elements) == 1) {
			return false
		}

		for i, component := range elements {
			if IsHiddenBasename(component) {
				return false
			}
			notLast := i < len(elements)-1
			if notLast && IsSpecialFolder(component) {
				return false
			}
		}

		return true
	}
}

func NewMetaNecessityInUnityProjectSubDir() MetaNecessity {
	return func(targetPath typedpath.SlashPath) bool {
		elements := pathutil.SplitPathElements(targetPath)
		if len(elements) == 0 {
			return false
		}

		for i, component := range elements {
			if IsHiddenBasename(component) {
				return false
			}
			notLast := i < len(elements)-1
			if notLast && IsSpecialFolder(component) {
				return false
			}
		}

		return true
	}
}

// IsHiddenBasename return whether the specified basename is treated as "hidden" by Unity.
// > During the import process, Unity completely ignores the following files and folders in the Assets folder (or a sub-folder within it):
// >
// > Hidden folders.
// > Files and folders which start with ‘.’.
// > Files and folders which end with ‘~’.
// > Files and folders named cvs.
// > Files with the extension .tmp.
// https://docs.unity3d.com/2020.2/Documentation/Manual/SpecialFolders.html
func IsHiddenBasename(baseName typedpath.BaseName) bool {
	return IsMeta(typedpath.SlashPathFromBaseName(baseName)) ||
		strings.HasPrefix(string(baseName), ".") ||
		strings.HasSuffix(string(baseName), "~") ||
		strings.HasSuffix(string(baseName), ".tmp")
}

// IsSpecialFolder return whether the specified basename is treated as "special" by Unity.
// Unity requires .meta files for only the special folder itself, and not requires .meta files for contents of the special folder.
// https://docs.unity3d.com/ja/2023.1/Manual/SpecialFolders.html
// https://forum.unity.com/threads/loadable-plugin-directory-import-behaviour-change-androidlib-bundle-framework-and-plugin.1381113/
func IsSpecialFolder(baseName typedpath.BaseName) bool {
	return strings.HasSuffix(string(baseName), ".xcframework") ||
		strings.HasSuffix(string(baseName), ".framework") ||
		strings.HasSuffix(string(baseName), ".androidlib") ||
		strings.HasSuffix(string(baseName), ".androidpack") ||
		strings.HasSuffix(string(baseName), ".bundle") ||
		strings.HasSuffix(string(baseName), ".plugin")
}

const MetaSuffix string = ".meta"

func IsMeta(path typedpath.SlashPath) bool {
	return strings.HasSuffix(string(path), MetaSuffix)
}

func TrimMetaFromSlash(path typedpath.SlashPath) typedpath.SlashPath {
	return typedpath.SlashPath(strings.TrimSuffix(string(path), MetaSuffix))
}

func TrimMetaFromRaw(path typedpath.RawPath) typedpath.RawPath {
	return typedpath.RawPath(strings.TrimSuffix(string(path), MetaSuffix))
}

func MetaPath(path typedpath.SlashPath) typedpath.SlashPath {
	return typedpath.SlashPath(fmt.Sprintf("%s%s", path, MetaSuffix))
}
