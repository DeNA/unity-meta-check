package filecollector

import (
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/util/chanutil"
	"github.com/DeNA/unity-meta-check/util/logging"
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
`, nil)
	opts := &Options{IgnoreCase: false}
	onlyTracked := New(gitLsFiles, spyLogger)

	spyCh := make(chan Entry)
	actual := make([]Entry, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		actual = chanutil.ToSlice(spyCh)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Path < actual[j].Path
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

	expected := []Entry{
		// NOTE: Directory should be added automatically (git ls-files does not print any directories)
		{Path: "dir1", IsDir: true},
		{Path: "dir1/fileB", IsDir: false},
		{Path: "fileA", IsDir: false},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
