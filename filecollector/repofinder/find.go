package repofinder

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
	"path/filepath"
)

type RepositoryType bool

const (
	RepositoryTypeIsSubmodule RepositoryType = true
	RepositoryTypeIsNested    RepositoryType = false
)

type RepoFinder func() ([]FoundRepo, error)

func New(rootDirAbs typedpath.RawPath, targetDirRel typedpath.RawPath) RepoFinder {
	return func() ([]FoundRepo, error) {
		var targetDirAbs typedpath.RawPath
		if targetDirRel == "." {
			targetDirAbs = rootDirAbs
		} else {
			targetDirAbs = rootDirAbs.JoinRawPath(targetDirRel)
		}

		result := make([]FoundRepo, 0)

		if err := filepath.Walk(string(targetDirAbs), func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}

			// NOTE: This matches both submodules and nested repositories. The "nested repository" means a repository
			//       that is cloned into other repository with no any submodule commands. The nested repositories are
			//       known as anti-pattern, but it is frequently needed for game developers.
			if info.Name() == ".git" {
				relPath, err := rootDirAbs.Rel(typedpath.RawPath(path))
				if err != nil {
					return err
				}
				isDir := info.IsDir()

				// NOTE: Ignore the repository at the rootDirAbs.
				if relPath != ".git" {
					result = append(result, FoundRepo{
						Type:    RepositoryType(!isDir),
						RawPath: relPath.Dir(),
					})
				}
				// NOTE: Skip finding into .git if it is a directory (this is a nested repository). But if it is a file
				//       (this is a submodule) skipping will skip unexpectedly the parent directory instead of the file.
				if isDir {
					return filepath.SkipDir
				}
				return nil
			}
			return nil
		}); err != nil {
			return result, err
		}
		return result, nil
	}
}
