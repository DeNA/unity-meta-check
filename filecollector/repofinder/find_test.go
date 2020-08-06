package repofinder

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestFind(t *testing.T) {
	testDir := setUpTestDir()
	spyCh := make(chan *FoundRepo)

	var actual []*FoundRepo
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for found := range spyCh {
			actual = append(actual, found)
		}
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].RawPath < actual[j].RawPath
		})
		wg.Done()
	}()

	findNested := New(testDir, ".")
	err := findNested(spyCh)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}
	close(spyCh)

	wg.Wait()
	expected := []*FoundRepo{
		{RepositoryTypeIsNested, "nested1"},
		{RepositoryTypeIsNested, typedpath.NewRawPath("nested1", "nestedInNested1")},
		{RepositoryTypeIsSubmodule, "nested2"},
		{RepositoryTypeIsSubmodule, typedpath.NewRawPath("nested2", "nestedInNested2")},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}


func TestFindOnRel(t *testing.T) {
	testDir := setUpTestDir()
	spyCh := make(chan *FoundRepo)

	var actual []*FoundRepo
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for found := range spyCh {
			actual = append(actual, found)
		}
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].RawPath < actual[j].RawPath
		})
		wg.Done()
	}()

	findNested := New(testDir, "nested1")
	err := findNested(spyCh)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}
	close(spyCh)

	wg.Wait()
	expected := []*FoundRepo{
		{RepositoryTypeIsNested, "nested1"},
		{RepositoryTypeIsNested, typedpath.NewRawPath("nested1", "nestedInNested1")},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}

func setUpTestDir() typedpath.RawPath {
	workDir, err := ioutil.TempDir(os.TempDir(), "unity-meta-check-tests.")
	if err != nil {
		panic(err.Error())
	}

	mkdirp(workDir, ".git")
	mkdirp(workDir, "nested1")
	mkdirp(workDir, "nested1", ".git")
	mkdirp(workDir, "nested1", "nestedInNested1")
	mkdirp(workDir, "nested1", "nestedInNested1", ".git")
	mkdirp(workDir, "nested2")
	touch(workDir, "nested2", ".git") // Means submodule.
	mkdirp(workDir, "nested2", "nestedInNested2")
	touch(workDir, "nested2", "nestedInNested2", ".git") // Means submodule.
	mkdirp(workDir, "others")

	return typedpath.NewRawPathUnsafe(workDir)
}

func mkdirp(cwd string, path ...string) {
	path = append([]string{cwd}, path...)
	if err := os.Mkdir(filepath.Join(path...), 0755); err != nil {
		panic(err.Error())
	}
}

func touch(cwd string, path ...string) {
	path = append([]string{cwd}, path...)
	file, err := os.Create(filepath.Join(path...))
	if err != nil {
		panic(err.Error())
	}
	if err := file.Close(); err != nil {
		panic(err.Error())
	}
}
