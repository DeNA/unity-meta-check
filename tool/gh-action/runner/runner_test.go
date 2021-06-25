package runner

import (
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/junit"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"net/url"
	"testing"
)

func TestNewRunner(t *testing.T) {
	type Case struct {
		Options                        *Options
		CheckResult                    *checker.CheckResult
		ExpectedJUnitWriteToFileCalled bool
		ExpectedSendFuncCalled         bool
		ExpectedAutoFixerCalled        bool
		ExpectedOK                     bool
	}

	cases := map[string]Case{}

	for _, targetType := range []checker.TargetType{checker.TargetTypeIsUnityProjectRootDirectory, checker.TargetTypeIsUnityProjectSubDirectory} {
		for _, missingMeta := range [][]typedpath.SlashPath{{}, {"Assets/Missing.meta"}, {ignoredPath}} {
			for _, danglingMeta := range [][]typedpath.SlashPath{{}, {"Assets/Dangling.meta"}, {ignoredPath}} {
				for _, enableJUnit := range []bool{false, true} {
					for _, enablePRComment := range []bool{false, true} {
						for _, enableAutofix := range []bool{false, true} {
							desc := fmt.Sprintf(
								"missingMeta=%t, danglingMeta=%t, enableJUnit=%t, enablePRComment=%t, enableAutofix=%t",
								len(missingMeta) > 0,
								len(danglingMeta) > 0,
								enableJUnit,
								enablePRComment,
								enableAutofix,
							)

							rootDirAbs := typedpath.NewRawPath("/", "path", "to", "proj")

							var junitOutPath typedpath.RawPath
							if enableJUnit {
								junitOutPath = typedpath.NewRawPath("out", "junit.xml")
							}
							
							var prCommentOpts *github.Options
							if enablePRComment {
								prCommentOpts = &github.Options{
									Tmpl:          &l10n.En,
									SendIfSuccess: true,
									Token:         "T0K3N",
									APIEndpoint:   mustURLParse("https://api.github.com"),
									Owner:         "octocat",
									Repo:          "Hello-World",
									PullNumber:    123,
								}
							}

							var autofixOpts *autofix.Options
							if enableAutofix {
								autofixOpts = &autofix.Options{
									RootDirAbs:   rootDirAbs,
									RootDirRel:   ".",
									AllowedGlobs: []globs.Glob{"."},
								}
							}

							cases[desc] = Case{
								Options: &Options{
									RootDirAbs: rootDirAbs,
									CheckerOpts: &checker.Options{
										IgnoreCase:                false,
										IgnoreSubmodulesAndNested: false,
										TargetType:                targetType,
									},
									FilterOpts: &resultfilter.Options{
										IgnoreDangling: false,
										IgnoredGlobs:   []globs.Glob{ignoredGlob},
										IgnoreCase:     false,
									},
									EnableJUnit:     enableJUnit,
									JUnitOutPath:    junitOutPath,
									EnablePRComment: enablePRComment,
									PRCommentOpts:   prCommentOpts,
									EnableAutofix:   enableAutofix,
									AutofixOpts:     autofixOpts,
								},
								CheckResult: &checker.CheckResult{
									MissingMeta:  missingMeta,
									DanglingMeta: danglingMeta,
								},
								ExpectedJUnitWriteToFileCalled: enableJUnit,
								ExpectedSendFuncCalled:         enablePRComment,
								ExpectedAutoFixerCalled:        enableAutofix,
								ExpectedOK:                     len(dropIgnored(missingMeta)) == 0 && len(dropIgnored(danglingMeta)) == 0,
							}
						}
					}
				}
			}
		}
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			spyLogger := logging.SpyLogger()
			buf := &bytes.Buffer{}

			junitWriteToFileCallArgs := make([]junit.WriteToFileCallArgs, 0)
			sendFuncCallArgs := make([]github.SendFuncCallArgs, 0)
			autoFixerCallArgs := make([]autofix.AutoFixerCallArgs, 0)

			run := NewRunner(
				checker.StubChecker(c.CheckResult, nil),
				resultfilter.NewFilter(ostestable.NewGetwd(), spyLogger),
				junit.SpyWriteToFileFunc(junit.StubWriteToFileFunc(nil), &junitWriteToFileCallArgs),
				github.SpySendFunc(github.StubSendFunc(nil), &sendFuncCallArgs),
				autofix.SpyAutoFixer(autofix.StubAutoFixer(nil), &autoFixerCallArgs),
				buf,
			)

			ok, err := run(c.Options)
			if err != nil {
				t.Log(spyLogger.Logs.String())
				t.Errorf("want nil, got %#v", err)
				return
			}

			if len(junitWriteToFileCallArgs) > 0 != c.ExpectedJUnitWriteToFileCalled {
				t.Log(spyLogger.Logs.String())
				t.Log(junitWriteToFileCallArgs)
				t.Errorf("want len(junitWriteToFileCallArgs) > 0 == %t, got %t", c.ExpectedJUnitWriteToFileCalled, len(junitWriteToFileCallArgs) > 0)
			}

			if len(sendFuncCallArgs) > 0 != c.ExpectedSendFuncCalled {
				t.Log(spyLogger.Logs.String())
				t.Log(sendFuncCallArgs)
				t.Errorf("want len(sendFuncCallArgs) > 0 == %t, got %t", c.ExpectedSendFuncCalled, len(sendFuncCallArgs) > 0)
			}

			if len(autoFixerCallArgs) > 0 != c.ExpectedAutoFixerCalled {
				t.Log(spyLogger.Logs.String())
				t.Log(autoFixerCallArgs)
				t.Errorf("want len(autoFixerCallArgs) > 0 == %t, got %t", c.ExpectedAutoFixerCalled, len(autoFixerCallArgs) > 0)
			}

			if ok != c.ExpectedOK {
				t.Log(spyLogger.Logs.String())
				t.Errorf("want ok == %t, got %t", c.ExpectedOK, ok)
			}
		})
	}
}

const (
	ignoredPath typedpath.SlashPath = "ignored"
	ignoredGlob globs.Glob          = "ignored*"
)

func dropIgnored(paths []typedpath.SlashPath) []typedpath.SlashPath {
	result := make([]typedpath.SlashPath, 0, len(paths))
	for _, p := range paths {
		if p == ignoredPath {
			continue
		}
		result = append(result, p)
	}
	return result
}

func mustURLParse(unsafeURL string) github.APIEndpoint {
	u, err := url.Parse(unsafeURL)
	if err != nil {
		panic(err.Error())
	}
	return u
}