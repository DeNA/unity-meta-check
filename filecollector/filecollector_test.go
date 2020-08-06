package filecollector

import (
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

func TestNewOnlyTracked(t *testing.T) {
	spyLogger := logging.SpyLogger()

	// NOTE: git returns a slash separated path on Windows.
	gitLsFiles := git.StubRawLsFiles(`dir1/fileB
fileA
dir2
`, nil)
	opts := &Options{IgnoreCase: false}
	onlyTracked := New(gitLsFiles, spyLogger)

	spyCh := make(chan typedpath.SlashPath)
	actual := make([]typedpath.SlashPath, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		actual = pathchan.ToSlice(spyCh)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i] < actual[j]
		})
		defer wg.Done()
	}()

	if err := onlyTracked("/path/to/repo", ".", opts, spyCh); err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}
	close(spyCh)
	wg.Wait()

	expected := []typedpath.SlashPath{
		"dir1", // Directory should be added automatically (git ls-files does not print any directories)
		"dir1/fileB",
		"dir2",
		"fileA",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
