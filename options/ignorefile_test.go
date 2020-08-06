package options

import (
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestNewIgnoredPathsBuilderSpecifiedAndIgnoreFileExists(t *testing.T) {
	spyLogger := logging.SpyLogger()
	buildIgnorePaths := NewIgnoredPathsBuilder(spyLogger)

	actual, err := buildIgnorePaths(
		typedpath.NewRawPath("testdata", "ignorefile", "ProjectHasMetaCheckIgnore", ".meta-check-ignore"),
		"/path/to/any",
	)

	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := []globs.Glob{
		"path/to/ignored",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}

func TestNewIgnoredPathsBuilderSpecifiedButIgnoreFileDoesNotExist(t *testing.T) {
	spyLogger := logging.SpyLogger()
	buildIgnorePaths := NewIgnoredPathsBuilder(spyLogger)

	_, err := buildIgnorePaths(
		typedpath.NewRawPath("testdata", "ignorefile", "ProjectDoesNotHaveMetaCheckIgnore", ".meta-check-ignore"),
		"/path/to/any",
	)

	if err == nil {
		t.Log(spyLogger.Logs.String())
		t.Error("want error, got nil")
		return
	}
}

func TestNewIgnoredPathsBuilderOmitButIgnoreFileExists(t *testing.T) {
	spyLogger := logging.SpyLogger()
	buildIgnorePaths := NewIgnoredPathsBuilder(spyLogger)

	actual, err := buildIgnorePaths(
		"",
		typedpath.NewRawPath("testdata", "ignorefile", "ProjectHasMetaCheckIgnore"),
	)

	if err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := []globs.Glob{
		"path/to/ignored",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}

func TestNewIgnoredPathsBuilderOmitAndIgnoreFileDoesNotExist(t *testing.T) {
	spyLogger := logging.SpyLogger()
	buildIgnorePaths := NewIgnoredPathsBuilder(spyLogger)

	actual, err := buildIgnorePaths(
		"",
		typedpath.NewRawPath("testdata", "meta-check-ignore", "ProjectDoesNotHaveMetaCheckIgnore"),
	)

	if err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := make([]globs.Glob, 0)
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
