package options

import (
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestBuildOptions(t *testing.T) {
	githubComAPIEndpoint, err := url.Parse("https://api.github.com")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	githubEnterpriseServerAPIEndpoint, err := url.Parse("https://github.example.com/api/v3")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	cases := map[string]struct {
		Args     []string
		Env      map[string]string
		Expected *Options
	}{
		"neither -lang or -template-file": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityInfo,
				Tmpl:          &l10n.En,
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubComAPIEndpoint,
				SendIfSuccess: true,
			},
		},
		"-lang": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1", "-lang", "ja"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityInfo,
				Tmpl:          &l10n.Ja,
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubComAPIEndpoint,
				SendIfSuccess: true,
			},
		},
		"-template-file": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1", "-template-file", "testdata/example-template.json"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityInfo,
				Tmpl:          &l10n.Template{
					SuccessMessage: "SUCCESS",
					FailureMessage: "FAILURE",
					StatusHeader:   "HEADER_STATUS",
					FilePathHeader: "HEADER_FILE_PATH",
					StatusMissing:  "STATUS_MISSING",
					StatusDangling: "STATUS_DANGLING",
				},
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubComAPIEndpoint,
				SendIfSuccess: true,
			},
		},
		"-silent": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1", "-silent"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityWarn,
				Tmpl:          &l10n.En,
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubComAPIEndpoint,
				SendIfSuccess: true,
			},
		},
		"-debug": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1", "-debug"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityDebug,
				Tmpl:          &l10n.En,
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubComAPIEndpoint,
				SendIfSuccess: true,
			},
		},
		"both -silent and -debug": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1", "-debug", "-silent"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityDebug, // -debug win
				Tmpl:          &l10n.En,
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubComAPIEndpoint,
				SendIfSuccess: true,
			},
		},
		"-no-send-success": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1", "-no-send-success"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityInfo,
				Tmpl:          &l10n.En,
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubComAPIEndpoint,
				SendIfSuccess: false,
			},
		},
		"-api-endpoint": {
			Args: []string{"-owner", "octocat", "-repo", "Hello-World", "-pull", "1", "-api-endpoint", "https://github.example.com/api/v3"},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				LogLevel:      logging.SeverityInfo,
				Tmpl:          &l10n.En,
				Token:         "T0K3N",
				Owner:         "octocat",
				Repo:          "Hello-World",
				PullNumber:    1,
				APIEndpoint:   githubEnterpriseServerAPIEndpoint,
				SendIfSuccess: true,
			},
		},
		"-version": {
			Args: []string{"-version"},
			Expected: &Options{
				Version: true,
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			stdin := io.NopCloser(strings.NewReader(""))
			stdout := testutil.SpyWriteCloser(nil)
			stderr := testutil.SpyWriteCloser(nil)
			procInout := cli.ProcessInout{
				Stdin:  stdin,
				Stdout: stdout,
				Stderr: stderr,
			}

			opts, err := BuildOptions(c.Args, procInout, cli.StubEnv(c.Env))
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
