package filecollector

import (
	"github.com/DeNA/unity-meta-check/filecollector/repofinder"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/pathchan"
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
	findNested := repofinder.Const([]*repofinder.FoundRepo{
		{Type: repofinder.RepositoryTypeIsSubmodule, RawPath: "nested1"},
		{Type: repofinder.RepositoryTypeIsNested, RawPath: "nested2"},
		{Type: repofinder.RepositoryTypeIsNested, RawPath: typedpath.NewRawPath("nested1", "nestedInNested")},
	}, nil)
	spyLogger := logging.SpyLogger()
	collectRecursive := NewFileAggregator(gitLsFiles, findNested, spyLogger)

	var actual []typedpath.SlashPath
	spyCh := make(chan typedpath.SlashPath)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		actual = pathchan.ToSlice(spyCh)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i] < actual[j]
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

	expected := []typedpath.SlashPath{
		"dir",
		"dir/file2",
		"file1",
		"nested1",
		"nested1/nested1-file",
		"nested1/nestedInNested",
		"nested1/nestedInNested/nestedInNested-file",
		"nested2",
		"nested2/nested2-file",
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
	findNested := repofinder.Empty
	spyLogger := logging.SpyLogger()
	collectRecursive := NewFileAggregator(gitLsFiles, findNested, spyLogger)

	var actual []typedpath.SlashPath
	spyCh := make(chan typedpath.SlashPath)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		actual = pathchan.ToSlice(spyCh)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i] < actual[j]
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

	expected := []typedpath.SlashPath{
		"path",
		"path/to",
		"path/to/file1",
		"path/to/file2",
		"path/to/file3",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
