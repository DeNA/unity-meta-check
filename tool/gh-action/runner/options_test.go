package runner

import (
	"github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"net/url"
	"reflect"
	"testing"
)

func TestNewValidateFunc(t *testing.T) {
	githubComAPIEndpoint, err := url.Parse("https://api.github.com")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	tmpl := &l10n.Template{
		SuccessMessage: "s",
		FailureMessage: "f",
		StatusHeader:   "s",
		FilePathHeader: "f",
		StatusMissing:  "s",
		StatusDangling: "s",
	}

	cases := map[string]struct {
		RootDirAbs         typedpath.SlashPath
		RootDirRel         typedpath.SlashPath
		Inputs             inputs.Inputs
		Token              github.Token
		DetectedTargetType checker.TargetType
		BuiltIgnoredGlobs  []globs.Glob
		ReadTmpl           *l10n.Template
		Expected           *Options
	}{
		"all default": {
			RootDirAbs: "/path/to/project",
			Inputs: inputs.Inputs{
				LogLevel:   "INFO",
				TargetType: "auto-detect",
			},
			DetectedTargetType: checker.TargetTypeIsUnityProjectRootDirectory,
			BuiltIgnoredGlobs:  []globs.Glob{"ignore*"},
			Expected: &Options{
				RootDirAbs: typedpath.NewRootRawPath("path", "to", "project"),
				CheckerOpts: &checker.Options{
					IgnoreCase:                false,
					IgnoreSubmodulesAndNested: false,
					// NOTE: same as the value of DetectedTargetType.
					TargetType: checker.TargetTypeIsUnityProjectRootDirectory,
				},
				FilterOpts: &resultfilter.Options{
					IgnoreDangling: false,
					// NOTE: same as the value of BuiltIgnoredGlobs.
					IgnoredGlobs: []globs.Glob{"ignore*"},
					IgnoreCase:   false,
				},
				EnableJUnit:     false,
				EnablePRComment: false,
				EnableAutofix:   false,
			},
		},
		"ignore-case": {
			RootDirAbs: "/path/to/project",
			Inputs: inputs.Inputs{
				LogLevel:   "INFO",
				TargetType: "auto-detect",
				IgnoreCase: true,
			},
			DetectedTargetType: checker.TargetTypeIsUnityProjectRootDirectory,
			BuiltIgnoredGlobs:  []globs.Glob{},
			Expected: &Options{
				RootDirAbs: typedpath.NewRootRawPath("path", "to", "project"),
				CheckerOpts: &checker.Options{
					IgnoreCase:                true,
					IgnoreSubmodulesAndNested: false,
					// NOTE: same as the value of DetectedTargetType.
					TargetType: checker.TargetTypeIsUnityProjectRootDirectory,
				},
				FilterOpts: &resultfilter.Options{
					IgnoreDangling: false,
					// NOTE: same as the value of BuiltIgnoredGlobs.
					IgnoredGlobs: []globs.Glob{},
					IgnoreCase:   true,
				},
				EnableJUnit:     false,
				EnablePRComment: false,
				EnableAutofix:   false,
			},
		},
		"enable junit": {
			RootDirAbs: "/path/to/project",
			Inputs: inputs.Inputs{
				LogLevel:     "INFO",
				TargetType:   "auto-detect",
				EnableJUnit:  true,
				JUnitXMLPath: "/path/to/project/junit.xml",
			},
			DetectedTargetType: checker.TargetTypeIsUnityProjectRootDirectory,
			BuiltIgnoredGlobs:  []globs.Glob{},
			Expected: &Options{
				RootDirAbs: typedpath.NewRootRawPath("path", "to", "project"),
				CheckerOpts: &checker.Options{
					IgnoreCase:                false,
					IgnoreSubmodulesAndNested: false,
					// NOTE: same as the value of DetectedTargetType.
					TargetType: checker.TargetTypeIsUnityProjectRootDirectory,
				},
				FilterOpts: &resultfilter.Options{
					IgnoreDangling: false,
					// NOTE: same as the value of BuiltIgnoredGlobs.
					IgnoredGlobs: []globs.Glob{},
					IgnoreCase:   false,
				},
				EnableJUnit:     true,
				JUnitOutPath:    "/path/to/project/junit.xml",
				EnablePRComment: false,
				EnableAutofix:   false,
			},
		},
		"enable pr-comment with lang": {
			RootDirAbs: "/path/to/project",
			Inputs: inputs.Inputs{
				LogLevel:             "INFO",
				TargetType:           "auto-detect",
				EnablePRComment:      true,
				PRCommentPRNumber:    123,
				PRCommentLang:        "en",
				PRCommentOwner:       "octocat",
				PRCommentRepo:        "Hello-World",
				PRCommentAPIEndpoint: "https://api.github.com",
				PRCommentSendSuccess: true,
			},
			Token:              "T0K3N",
			DetectedTargetType: checker.TargetTypeIsUnityProjectRootDirectory,
			BuiltIgnoredGlobs:  []globs.Glob{},
			Expected: &Options{
				RootDirAbs: typedpath.NewRootRawPath("path", "to", "project"),
				CheckerOpts: &checker.Options{
					IgnoreCase:                false,
					IgnoreSubmodulesAndNested: false,
					// NOTE: same as the value of DetectedTargetType.
					TargetType: checker.TargetTypeIsUnityProjectRootDirectory,
				},
				FilterOpts: &resultfilter.Options{
					IgnoreDangling: false,
					// NOTE: same as the value of BuiltIgnoredGlobs.
					IgnoredGlobs: []globs.Glob{},
					IgnoreCase:   false,
				},
				EnableJUnit:     false,
				EnablePRComment: true,
				PRCommentOpts: &github.Options{
					Tmpl:          &l10n.En,
					SendIfSuccess: true,
					Token:         "T0K3N",
					APIEndpoint:   githubComAPIEndpoint,
					Owner:         "octocat",
					Repo:          "Hello-World",
					PullNumber:    123,
				},
			},
		},
		"enable pr-comment with a template file": {
			RootDirAbs: "/path/to/project",
			Inputs: inputs.Inputs{
				LogLevel:              "INFO",
				TargetType:            "auto-detect",
				EnablePRComment:       true,
				PRCommentPRNumber:     123,
				PRCommentTmplFilePath: "path/to/tmpl.json",
				PRCommentOwner:        "octocat",
				PRCommentRepo:         "Hello-World",
				PRCommentAPIEndpoint:  "https://api.github.com",
				PRCommentSendSuccess:  true,
			},
			Token:              "T0K3N",
			DetectedTargetType: checker.TargetTypeIsUnityProjectRootDirectory,
			BuiltIgnoredGlobs:  []globs.Glob{},
			ReadTmpl: tmpl,
			Expected: &Options{
				RootDirAbs: typedpath.NewRootRawPath("path", "to", "project"),
				CheckerOpts: &checker.Options{
					IgnoreCase:                false,
					IgnoreSubmodulesAndNested: false,
					// NOTE: same as the value of DetectedTargetType.
					TargetType: checker.TargetTypeIsUnityProjectRootDirectory,
				},
				FilterOpts: &resultfilter.Options{
					IgnoreDangling: false,
					// NOTE: same as the value of BuiltIgnoredGlobs.
					IgnoredGlobs: []globs.Glob{},
					IgnoreCase:   false,
				},
				EnableJUnit:     false,
				EnablePRComment: true,
				PRCommentOpts: &github.Options{
					Tmpl:          tmpl,
					SendIfSuccess: true,
					Token:         "T0K3N",
					APIEndpoint:   githubComAPIEndpoint,
					Owner:         "octocat",
					Repo:          "Hello-World",
					PullNumber:    123,
				},
			},
		},
		"enable autofix": {
			RootDirAbs: "/path/to/project",
			RootDirRel: ".",
			Inputs: inputs.Inputs{
				LogLevel:      "INFO",
				TargetType:    "auto-detect",
				EnableAutofix: true,
				AutofixGlobs:  []string{"Assets/Allow/To/Fix/*"},
			},
			DetectedTargetType: checker.TargetTypeIsUnityProjectRootDirectory,
			BuiltIgnoredGlobs:  []globs.Glob{},
			Expected: &Options{
				RootDirAbs: typedpath.NewRootRawPath("path", "to", "project"),
				CheckerOpts: &checker.Options{
					IgnoreCase:                false,
					IgnoreSubmodulesAndNested: false,
					// NOTE: same as the value of DetectedTargetType.
					TargetType: checker.TargetTypeIsUnityProjectRootDirectory,
				},
				FilterOpts: &resultfilter.Options{
					IgnoreDangling: false,
					// NOTE: same as the value of BuiltIgnoredGlobs.
					IgnoredGlobs: []globs.Glob{},
					IgnoreCase:   false,
				},
				EnableJUnit:     false,
				EnablePRComment: false,
				EnableAutofix:   true,
				AutofixOpts: &autofix.Options{
					RootDirAbs:   "/path/to/project",
					RootDirRel:   ".",
					AllowedGlobs: []globs.Glob{"Assets/Allow/To/Fix/*"},
				},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			validate := NewValidateFunc(
				options.StubUnityProjectDetector(c.DetectedTargetType, nil),
				options.StubIgnoredPathBuilder(c.BuiltIgnoredGlobs, nil),
				autofix.StubOptionsBuilderWithRootDirAbsAndRel(c.RootDirRel.ToRaw()),
				l10n.StubTemplateFileReader(c.ReadTmpl, nil),
			)

			opts, err := validate(c.RootDirAbs.ToRaw(), c.Inputs, c.Token)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if !reflect.DeepEqual(opts, c.Expected) {
				t.Error(cmp.Diff(c.Expected, opts))
			}
		})
	}
}
