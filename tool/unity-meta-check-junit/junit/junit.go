package junit

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/jstemmer/go-junit-report/formatter"
	"github.com/jstemmer/go-junit-report/parser"
	"io"
	"os"
	"runtime"
	"time"
)

func WriteToFile(result *checker.CheckResult, startTime time.Time, outPath typedpath.RawPath) error {
	if err := os.MkdirAll(string(outPath.Dir()), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(string(outPath), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(){ _ = file.Close() }()

	endTime := time.Now()
	return Write(result, endTime.Sub(startTime), file)
}

func Write(result *checker.CheckResult, duration time.Duration, writer io.Writer) error {
	maxLen := result.Len()
	var packages []parser.Package
	if maxLen == 0 {
		packages = []parser.Package{
			{
				Name: "unity-meta-check",
				Tests: []*parser.Test{
					{
						Name:   "OK",
						Result: parser.PASS,
						Output: []string{"No missing or dangling .meta exist. Perfect!"},
						Duration: duration,
					},
				},
			},
		}
	} else {
		durationAvg := time.Duration(int(duration) / maxLen)
		packages = make([]parser.Package, maxLen)
		i := 0
		for _, missingMeta := range result.MissingMeta {
			packages[i] = parser.Package{
				Name: string(unity.TrimMetaFromSlash(missingMeta)),
				Tests: []*parser.Test{
					{
						Name:   "meta",
						Result: parser.FAIL,
						Output: []string{
							fmt.Sprintf("File or directory exists: %s", unity.TrimMetaFromSlash(missingMeta)),
							fmt.Sprintf("But .meta is missing: %s", missingMeta),
						},
						Duration: durationAvg,
					},
				},
			}
			i++
		}
		for _, danglingMeta := range result.DanglingMeta {
			packages[i] = parser.Package{
				Name: string(unity.TrimMetaFromSlash(danglingMeta)),
				Tests: []*parser.Test{
					{
						Name:   "meta",
						Result: parser.FAIL,
						Output: []string{
							fmt.Sprintf("File or directory does not exist: %s", unity.TrimMetaFromSlash(danglingMeta)),
							fmt.Sprintf("But .meta is present: %s", danglingMeta),
						},
						Duration: durationAvg,
					},
				},
			}
			i++
		}
	}

	junitReport := &parser.Report{Packages: packages}

	if err := formatter.JUnitReportXML(junitReport, false, runtime.Version(), writer); err != nil {
		return err
	}
	return nil
}
