package checker

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/filecollector"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/pathutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"sort"
)

type CheckingWorker func(rootDirAbs typedpath.RawPath, ignoreCase bool, reader <-chan filecollector.Entry) (*CheckResult, error)

func NewCheckingWorker(requiresMeta unity.MetaNecessity, logger logging.Logger) CheckingWorker {
	return func(rootDirAbs typedpath.RawPath, ignoreCase bool, reader <-chan filecollector.Entry) (*CheckResult, error) {
		entries, err := labelMeta(requiresMeta, reader, logger)
		if err != nil {
			return nil, err
		}

		expectedMetaSet := pathutil.NewPathSet(ignoreCase)
		allowedMetaSet := pathutil.NewPathSet(ignoreCase)
		actualMetaSet := pathutil.NewPathSet(ignoreCase)

		for _, entry := range entries {
			if entry.IsMeta {
				actualMetaSet.Add(entry.Path)
			} else if entry.NeedsMeta {
				expectedMetaSet.Add(unity.MetaPath(entry.Path))
				allowedMetaSet.Add(unity.MetaPath(entry.Path))
			} else {
				allowedMetaSet.Add(unity.MetaPath(entry.Path))
			}
		}

		return NewCheckResult(
			detectMissing(expectedMetaSet, actualMetaSet),
			detectDangling(allowedMetaSet, actualMetaSet, entries),
		), nil
	}
}

type entry struct {
	Path      typedpath.SlashPath
	IsMeta    bool
	NeedsMeta bool
	IsDir     bool
}

func labelMeta(requiresMetaFunc unity.MetaNecessity, reader <-chan filecollector.Entry, logger logging.Logger) ([]entry, error) {
	entries := make([]entry, 0)
	for targetEntry := range reader {
		isMeta := unity.IsMeta(targetEntry.Path)
		requiresMeta := requiresMetaFunc(targetEntry.Path)

		if isMeta {
			logger.Debug(fmt.Sprintf("meta found ... %s", targetEntry.Path))
		} else if requiresMeta {
			logger.Debug(fmt.Sprintf("needs meta ... %s", targetEntry.Path))
		} else {
			logger.Debug(fmt.Sprintf("skipped ... %s", targetEntry.Path))
		}

		entries = append(entries, entry{
			Path:      targetEntry.Path,
			IsMeta:    isMeta,
			NeedsMeta: requiresMeta,
			IsDir:     targetEntry.IsDir,
		})
	}
	return entries, nil
}

func detectMissing(expectedMetaSet *pathutil.PathSet, actualMetaSet *pathutil.PathSet) []typedpath.SlashPath {
	missingMeta := expectedMetaSet.Difference(actualMetaSet)
	sort.Slice(missingMeta, func(i, j int) bool {
		return missingMeta[i] < missingMeta[j]
	})
	return missingMeta
}

func detectDangling(allowedMetaSet *pathutil.PathSet, actualMetaSet *pathutil.PathSet, entries []entry) []typedpath.SlashPath {
	pairs := make([]pathutil.PathPair[entry], len(entries))
	for i, e := range entries {
		pairs[i] = pathutil.PathPair[entry]{
			Path:  e.Path,
			Value: e,
		}
	}

	danglingMeta := actualMetaSet.Difference(allowedMetaSet)

	// NOTE: Add directories that contains only dangling meta files.
	tree := pathutil.NewPathTreeWithValues(pairs...)
	nonDanglingDirs := make(map[typedpath.SlashPath]bool)
	_ = tree.Postorder(func(path typedpath.SlashPath, p pathutil.PathTreeEntry[entry]) error {
		// XXX: p.Value == nil indicates the path is complemented. The path get complemented is a directory always.
		isComplementedDir := p.Value == nil
		isDir := isComplementedDir || p.Value.IsDir
		if !isDir {
			return nil
		}

		hasNonDangling := false
		for baseName, subEntry := range p.Subtree {
			if subEntry.Value.IsDir {
				subEntryPath := path.JoinBaseName(baseName)
				// NOTE: nonDanglingDirs[subEntryPath] always exist because it is processed by postorder.
				if nonDanglingDirs[subEntryPath] {
					hasNonDangling = true
					break
				}
			} else {
				if !subEntry.Value.IsMeta {
					hasNonDangling = true
					break
				}
			}
		}

		if !hasNonDangling {
			meta := unity.MetaPath(path)
			if actualMetaSet.Has(meta) {
				danglingMeta = append(danglingMeta, meta)
			}
		}
		nonDanglingDirs[path] = hasNonDangling
		return nil
	})

	sort.Slice(danglingMeta, func(i, j int) bool {
		return danglingMeta[i] < danglingMeta[j]
	})

	return danglingMeta
}
