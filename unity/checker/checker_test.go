package checker

import (
	"github.com/DeNA/unity-meta-check/filecollector"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestCheck(t *testing.T) {
	strategy := Strategy{
		CollectFiles: filecollector.StubSuccessfulFileAggregator([]filecollector.Entry{
			{Path: "Assets/MissingMeta", IsDir: false},
			{Path: "Assets/DanglingMeta.meta", IsDir: false},
		}),
		RequiresMeta: unity.ConstMetaNecessity(true),
	}
	spyLogger := logging.SpyLogger()
	opts := &Options{
		IgnoreCase: false,
		TargetType: TargetTypeIsUnityProjectRootDirectory,
	}

	checker := newCheckerByStrategy(strategy, spyLogger)
	actual, err := checker("/path/to/Project", opts)

	if err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %s", err.Error())
		return
	}

	expected := &CheckResult{
		MissingMeta: []typedpath.SlashPath{
			"Assets/MissingMeta.meta",
		},
		DanglingMeta: []typedpath.SlashPath{
			"Assets/DanglingMeta.meta",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
