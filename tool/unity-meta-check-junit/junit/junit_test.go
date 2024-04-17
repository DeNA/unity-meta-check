package junit

import (
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"runtime"
	"testing"
)

func TestWrite(t *testing.T) {
	cases := map[string]struct {
		Result   *checker.CheckResult
		Expected string
	}{
		"empty (boundary)": {
			Result: checker.NewCheckResult([]typedpath.SlashPath{}, []typedpath.SlashPath{}),
			Expected: fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
	<testsuite tests="1" failures="0" time="0.000" name="unity-meta-check">
		<properties>
			<property name="go.version" value="%s"></property>
		</properties>
		<testcase classname="unity-meta-check" name="OK" time="0.000"></testcase>
	</testsuite>
</testsuites>
`, runtime.Version()),
		},
		"both missing and dangling (easy to test)": {
			Result: checker.NewCheckResult(
				[]typedpath.SlashPath{
					typedpath.NewSlashPathUnsafe("path/to/missing.meta"),
				},
				[]typedpath.SlashPath{
					typedpath.NewSlashPathUnsafe("path/to/dangling.meta"),
				},
			),
			Expected: fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
	<testsuite tests="1" failures="1" time="0.000" name="path/to/missing" file="path/to/missing">
		<properties>
			<property name="go.version" value="%s"></property>
		</properties>
		<testcase classname="missing" name="meta" time="0.000" file="path/to/missing">
			<failure message="Failed" type="">File or directory exists: path/to/missing&#xA;But .meta is missing: path/to/missing.meta</failure>
		</testcase>
	</testsuite>
	<testsuite tests="1" failures="1" time="0.000" name="path/to/dangling" file="path/to/dangling">
		<properties>
			<property name="go.version" value="%s"></property>
		</properties>
		<testcase classname="dangling" name="meta" time="0.000" file="path/to/dangling">
			<failure message="Failed" type="">File or directory does not exist: path/to/dangling&#xA;But .meta is present: path/to/dangling.meta</failure>
		</testcase>
	</testsuite>
</testsuites>
`, runtime.Version(), runtime.Version()),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			buf := &bytes.Buffer{}
			if err := Write(c.Result, buf); err != nil {
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
