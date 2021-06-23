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
		CollectFiles: filecollector.StubSuccessfulFileAggregator([]typedpath.SlashPath{
			"/path/to/Project/Assets/MissingMeta",
			"/path/to/Project/Assets/DanglingMeta.meta",
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
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := &CheckResult{
		MissingMeta: []typedpath.SlashPath{
			"/path/to/Project/Assets/MissingMeta.meta",
		},
		DanglingMeta: []typedpath.SlashPath{
			"/path/to/Project/Assets/DanglingMeta.meta",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Log(spyLogger.Logs.String())
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
