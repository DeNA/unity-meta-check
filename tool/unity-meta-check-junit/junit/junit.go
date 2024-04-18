package junit

import (
	"encoding/xml"
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io"
	"os"
	"runtime"
)

type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Tests      int        `xml:"tests,attr"`
	Failures   int        `xml:"failures,attr"`
	Time       string     `xml:"time,attr"`
	Name       string     `xml:"name,attr"`
	File       *string    `xml:"file,attr,omitempty"`
	Properties Properties `xml:"properties"`
	TestCases  []TestCase `xml:"testcase"`
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	ClassName string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      string   `xml:"time,attr"`
	File      *string  `xml:"file,attr,omitempty"`
	Failure   *Failure `xml:"failure,omitempty"`
}

type Failure struct {
	XMLName  xml.Name `xml:"failure"`
	Message  string   `xml:"message,attr"`
	Type     string   `xml:"type,attr"`
	Contents string   `xml:",chardata"`
}

type Properties struct {
	XMLName    xml.Name   `xml:"properties"`
	Properties []Property `xml:"property"`
}

type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

type WriteToFileFunc func(result *checker.CheckResult, outPath typedpath.RawPath) error

func WriteToFile(result *checker.CheckResult, outPath typedpath.RawPath) error {
	if err := os.MkdirAll(string(outPath.Dir()), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(string(outPath), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return Write(result, file)
}

func Write(result *checker.CheckResult, writer io.Writer) error {
	maxLen := result.Len()
	props := &Properties{
		Properties: []Property{
			{
				Name:  "go.version",
				Value: runtime.Version(),
			},
		},
	}
	var testSuites *TestSuites
	if maxLen == 0 {
		testSuites = &TestSuites{
			TestSuites: []TestSuite{
				{
					Name:       "unity-meta-check",
					Tests:      1,
					Failures:   0,
					Time:       "0.000",
					Properties: *props,
					TestCases: []TestCase{
						{
							ClassName: "unity-meta-check",
							Name:      "OK",
							Time:      "0.000",
						},
					},
				},
			},
		}
	} else {
		suites := make([]TestSuite, maxLen)
		i := 0
		for _, missingMeta := range result.MissingMeta {
			file := string(unity.TrimMetaFromSlash(missingMeta))
			suites[i] = TestSuite{
				Name:       file,
				Tests:      1,
				Time:       "0.000",
				Failures:   1,
				File:       &file,
				Properties: *props,
				TestCases: []TestCase{
					{
						ClassName: "missing",
						Name:      "meta",
						Time:      "0.000",
						File:      &file,
						Failure: &Failure{
							Message:  "Failed",
							Contents: fmt.Sprintf("File or directory exists: %s\nBut .meta is missing: %s", file, missingMeta),
						},
					},
				},
			}
			i++
		}
		for _, danglingMeta := range result.DanglingMeta {
			file := string(unity.TrimMetaFromSlash(danglingMeta))
			suites[i] = TestSuite{
				Name:       file,
				Tests:      1,
				Time:       "0.000",
				Failures:   1,
				File:       &file,
				Properties: *props,
				TestCases: []TestCase{
					{
						ClassName: "dangling",
						Name:      "meta",
						Time:      "0.000",
						File:      &file,
						Failure: &Failure{
							Message:  "Failed",
							Contents: fmt.Sprintf("File or directory does not exist: %s\nBut .meta is present: %s", file, danglingMeta),
						},
					},
				},
			}
			i++
		}
		testSuites = &TestSuites{
			TestSuites: suites,
		}
	}

	bs, err := xml.MarshalIndent(testSuites, "", "\t")
	if err != nil {
		return err
	}
	if _, err := io.WriteString(writer, xml.Header); err != nil {
		return err
	}
	if _, err := writer.Write(bs); err != nil {
		return err
	}
	if _, err := writer.Write([]byte{'\n'}); err != nil {
		return err
	}
	return nil
}
