package resultfilter

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestFilterResult(t *testing.T) {
	cases := []struct {
		Result   *checker.CheckResult
		Opts     *Options
		Expected *checker.CheckResult
	}{
		{
			Result: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{},
			),
			Opts: &Options{
				IgnoreDangling: false,
				IgnoredGlobs:   []globs.Glob{},
			},
			Expected: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{},
			),
		},

		{
			Result: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Not/Added.meta"},
				[]typedpath.SlashPath{},
			),
			Opts: &Options{
				IgnoreDangling: false,
				IgnoredGlobs:   []globs.Glob{},
			},
			Expected: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Not/Added.meta"},
				[]typedpath.SlashPath{},
			),
		},

		{
			Result: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Not/Added.meta"},
				[]typedpath.SlashPath{},
			),
			Opts: &Options{
				IgnoreDangling: false,
				IgnoredGlobs: []globs.Glob{
					"Assets/Not",
				},
			},
			Expected: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{},
			),
		},

		{
			Result: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Not/Added.meta"},
				[]typedpath.SlashPath{},
			),
			Opts: &Options{
				IgnoreDangling: false,
				IgnoredGlobs: []globs.Glob{
					"Assets/Not/Added.meta",
				},
			},
			Expected: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{},
			),
		},

		{
			Result: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{"Assets/Not/Added.meta"},
			),
			Opts: &Options{
				IgnoreDangling: true,
				IgnoredGlobs:   []globs.Glob{},
			},
			Expected: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{},
			),
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v %v -> %v", c.Result, c.Opts, c.Expected), func(t *testing.T) {
			spyLogger := logging.SpyLogger()
			filter := NewFilter(spyLogger)
			actual, err := filter(c.Result, c.Opts)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if !reflect.DeepEqual(actual, c.Expected) {
				t.Log(spyLogger.Logs.String())
				t.Error(cmp.Diff(c.Expected, actual))
				return
			}
		})
	}
}
