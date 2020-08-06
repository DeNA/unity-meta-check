package markdown

import (
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestWriteMarkdown(t *testing.T) {
	cases := []struct{
		Result *checker.CheckResult
		Expected string
	}{
		{
			Result: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
			Expected: "SUCCESS_MESSAGE\n",
		},
		{
			Result: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{
					"path/to/missing.meta",
				},
				DanglingMeta: []typedpath.SlashPath{},
			},
			Expected: strings.Join([]string{
				"FAILURE_MESSAGE",
				"",
				"| HEADER_STATUS | HEADER_FILE_PATH |",
				"|:--|:--|",
				"| STATUS_MISSING | `path/to/missing.meta` |",
				"",
			}, "\n"),
		},
		{
			Result: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{
					"path/to/dangling.meta",
				},
			},
			Expected: strings.Join([]string{
				"FAILURE_MESSAGE",
				"",
				"| HEADER_STATUS | HEADER_FILE_PATH |",
				"|:--|:--|",
				"| STATUS_DANGLING | `path/to/dangling.meta` |",
				"",
			}, "\n"),
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c), func(t *testing.T) {
			buf := &bytes.Buffer{}

			err := WriteMarkdown(c.Result, &l10n.Template{
				SuccessMessage: "SUCCESS_MESSAGE",
				FailureMessage: "FAILURE_MESSAGE",
				StatusHeader:   "HEADER_STATUS",
				FilePathHeader: "HEADER_FILE_PATH",
				StatusMissing:  "STATUS_MISSING",
				StatusDangling: "STATUS_DANGLING",
			}, buf)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			actual := buf.String()
			if actual != c.Expected {
				t.Error(cmp.Diff(c.Expected, actual))
				return
			}
		})
	}
}
