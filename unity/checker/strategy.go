package checker

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/filecollector"
	"github.com/DeNA/unity-meta-check/filecollector/repofinder"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type Strategy struct {
	CollectFiles filecollector.FileAggregator
	RequiresMeta unity.MetaNecessity
}

type StrategySelector func(rootDirAbs typedpath.RawPath, opts *Options) (Strategy, error)

// NewStrategySelector returns a checker strategy for either Unity projects or UPM packages.
func NewStrategySelector(findPackages unity.FindPackages, lsFiles git.LsFiles, logger logging.Logger) StrategySelector {
	return func(rootDirAbs typedpath.RawPath, opts *Options) (Strategy, error) {
		switch opts.TargetType {
		case TargetTypeIsUnityProjectRootDirectory:
			foundPackages, err := findPackages(rootDirAbs)
			if err != nil {
				return Strategy{}, fmt.Errorf("cannot find local packages: %#v", err)
			}

			findRepo := NewRepoFinderForUnityProj(rootDirAbs, opts, foundPackages)
			return Strategy{
				CollectFiles: filecollector.NewFileAggregator(lsFiles, findRepo, logger),
				RequiresMeta: unity.NewMetaNecessityInUnityProject(unity.FoundPackagesToSlashRelPaths(foundPackages)),
			}, nil

		case TargetTypeIsUnityProjectSubDirectory:
			findRepo := NewRepoFinderFactoryForUPM(rootDirAbs, opts)
			return Strategy{
				CollectFiles: filecollector.NewFileAggregator(lsFiles, findRepo, logger),
				RequiresMeta: unity.NewMetaNecessityInUnityProjectSubDir(),
			}, nil

		default:
			return Strategy{}, fmt.Errorf("unsupported checking type: %q", opts.TargetType)
		}
	}
}

func NewRepoFinderForUnityProj(rootDirAbs typedpath.RawPath, opts *Options, foundPackages []*unity.FoundPackage) repofinder.RepoFinder {
	if opts.IgnoreSubmodulesAndNested {
		return repofinder.StubRepoFinder(nil, nil)
	}

	repoFinders := make([]repofinder.RepoFinder, len(foundPackages)+1)
	repoFinders[0] = repofinder.New(rootDirAbs, typedpath.RawPath(unity.AssetsDirBaseName))
	i := 1
	for _, foundPkg := range foundPackages {
		repoFinders[i] = repofinder.New(rootDirAbs, foundPkg.RelPath)
		i++
	}
	findRepo := repofinder.Compose(repoFinders)
	return findRepo
}

func NewRepoFinderFactoryForUPM(rootDirAbs typedpath.RawPath, opts *Options) repofinder.RepoFinder {
	if opts.IgnoreSubmodulesAndNested {
		return repofinder.StubRepoFinder(nil, nil)
	}

	return repofinder.New(rootDirAbs, ".")
}
