package report

import (
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestNewLineReporterParser(t *testing.T) {
	cases := []*checker.CheckResult{
		checker.NewCheckResult([]typedpath.SlashPath{}, []typedpath.SlashPath{}),
		checker.NewCheckResult([]typedpath.SlashPath{
			"path/to/missing1.meta",
			"path/to/missing2.meta",
		}, []typedpath.SlashPath{}),
		checker.NewCheckResult([]typedpath.SlashPath{}, []typedpath.SlashPath{
			"path/to/dangling1.meta",
			"path/to/dangling2.meta",
		}),
		checker.NewCheckResult(
			[]typedpath.SlashPath{
				"path/to/missing1.meta",
				"path/to/missing2.meta",
			},
			[]typedpath.SlashPath{
				"path/to/dangling1.meta",
				"path/to/dangling2.meta",
			},
		),
	}

	for _, result := range cases {
		t.Run(fmt.Sprintf("%v", result), func(t *testing.T) {
			buf := &bytes.Buffer{}

			if err := WriteResult(buf, result); err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			parse := NewParser()
			actual := parse(buf)
			expected := result

			if !reflect.DeepEqual(actual, expected) {
				t.Error(cmp.Diff(expected, actual))
				return
			}
		})
	}
}
