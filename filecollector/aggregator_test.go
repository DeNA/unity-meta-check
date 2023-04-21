package filecollector

import (
	"github.com/DeNA/unity-meta-check/filecollector/repofinder"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/util/chanutil"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestNewFileAggregator(t *testing.T) {
	gitLsFiles := git.FakeLsFiles(func(repoDir typedpath.RawPath) ([]string, error) {
		switch repoDir {
		case typedpath.NewSlashPathUnsafe("/path/to/repo").ToRaw():
			return []string{"file1", "dir/file2"}, nil
		case typedpath.NewSlashPathUnsafe("/path/to/repo/nested1").ToRaw():
			return []string{"nested1-file"}, nil
		case typedpath.NewSlashPathUnsafe("/path/to/repo/nested2").ToRaw():
			return []string{"nested2-file"}, nil
		case typedpath.NewSlashPathUnsafe("/path/to/repo/nested1/nestedInNested").ToRaw():
			return []string{"nestedInNested-file"}, nil
		default:
			panic(repoDir)
		}
	})
	findNested := repofinder.StubRepoFinder([]repofinder.FoundRepo{
		{Type: repofinder.RepositoryTypeIsSubmodule, RawPath: "nested1"},
		{Type: repofinder.RepositoryTypeIsNested, RawPath: "nested2"},
		{Type: repofinder.RepositoryTypeIsNested, RawPath: typedpath.NewRawPath("nested1", "nestedInNested")},
	}, nil)
	spyLogger := logging.SpyLogger()
	collectRecursive := NewFileAggregator(gitLsFiles, findNested, spyLogger)

	var actual []Entry
	spyCh := make(chan Entry)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		actual = chanutil.ToSlice(spyCh)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Path < actual[j].Path
		})
	}()

	projRootAbs := typedpath.NewSlashPathUnsafe("/path/to/repo").ToRaw()
	if err := collectRecursive(projRootAbs, &Options{IgnoreCase: false}, spyCh); err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}
	close(spyCh)

	wg.Wait()

	expected := []Entry{
		{Path: "dir", IsDir: true},
		{Path: "dir/file2", IsDir: false},
		{Path: "file1", IsDir: false},
		{Path: "nested1", IsDir: true},
		{Path: "nested1/nested1-file", IsDir: false},
		{Path: "nested1/nestedInNested", IsDir: true},
		{Path: "nested1/nestedInNested/nestedInNested-file", IsDir: false},
		{Path: "nested2", IsDir: true},
		{Path: "nested2/nested2-file", IsDir: false},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}

func TestNewFileAggregatorEmpty(t *testing.T) {
	gitLsFiles := git.StubLsFiles([]string{
		"path/to/file1",
		"path/to/file2",
		"path/to/file3",
	}, nil)
	findNested := repofinder.StubRepoFinder(nil, nil)
	spyLogger := logging.SpyLogger()
	collectRecursive := NewFileAggregator(gitLsFiles, findNested, spyLogger)

	var actual []Entry
	spyCh := make(chan Entry)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		actual = chanutil.ToSlice(spyCh)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Path < actual[j].Path
		})
	}()

	projRootAbs := typedpath.NewSlashPathUnsafe("/path/to/repo").ToRaw()
	if err := collectRecursive(projRootAbs, &Options{IgnoreCase: false}, spyCh); err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}
	close(spyCh)

	wg.Wait()

	expected := []Entry{
		{Path: "path", IsDir: true},
		{Path: "path/to", IsDir: true},
		{Path: "path/to/file1", IsDir: false},
		{Path: "path/to/file2", IsDir: false},
		{Path: "path/to/file3", IsDir: false},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
