package yaml

import (
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/version"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"testing"
)

func TestRecentActionYAML(t *testing.T) {
	actual, err := os.ReadFile("../testdata/action.yml")
	if err != nil {
		t.Errorf("want nil, got %v", err)
		return
	}

	buf := &bytes.Buffer{}
	if _, err := WriteTo(buf); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	if !reflect.DeepEqual(buf.Bytes(), actual) {
		t.Error(cmp.Diff(buf.String(), actual))
	}
}

// NOTE: This test is fragile, but we can use like Golden Testing.
func TestWriteInputsAsGHActionYAML(t *testing.T) {
	buf := &bytes.Buffer{}
	_, err := WriteInputsAsGHActionYAML(buf, BuildMetadata().Inputs)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := `
inputs:
  "log_level":
    description: "log level for unity-meta-check toolchains (available: ERROR/WARN/INFO/DEBUG)"
    default: "INFO"

  "target_type":
    description: "target type for unity-meta-check (available: auto-detect/unity-project/unity-project-sub-dir/upm-package)"
    default: "auto-detect"

  "ignore_dangling":
    description: "ignore dangling .meta"
    default: false

  "ignore_case":
    description: "do not treat case of file paths"
    default: false

  "ignore_submodules_and_nested":
    description: "ignore git submodules and nesting repositories (note: this is RECOMMENDED but not enabled by default because it can cause to miss problems in submodules or nesting repositories)"
    default: false

  "ignored_file_path":
    description: "path to .meta-check-ignore"
    default: ".meta-check-ignore"

  "enable_autofix":
    description: "enable autofix (note: it can repair very limited file types)"
    default: false

  "autofix_missing":
    description: "fix missing .meta"
    default: true

  "autofix_dangling":
    description: "fix dangling .meta"
    default: true

  "autofix_globs":
    description: "glob match directories to where to allow do autofix (note: required if enable_autofix is true)"
    default: ""

  "enable_junit":
    description: "enable JUnit-style reporting"
    default: false

  "junit_xml_path":
    description: "file path for generated JUnit test reports (note: required if enable_junit is true)"
    default: ""

  "enable_pr_comment":
    description: "enable reporting via GitHub Pull Request Comments"
    default: false

  "pr_comment_lang":
    description: "language code for GitHub Pull Request Comments (available: en/ja, note: cannot specify both lang and pr_comment_tmpl_file)"
    default: false

  "pr_comment_tmpl_file":
    description: "file path to custom template file for GitHub Pull Request Comments (note: cannot specify both lang and pr_comment_tmpl_file)"
    default: ""

  "pr_comment_owner":
    description: "owner of the GitHub repository (note: required if enable_pr_comment is true)"
    default: ""

  "pr_comment_repo":
    description: "name of the GitHub repository (note: required if enable_pr_comment is true)"
    default: ""

  "pr_comment_pull":
    description: "pull request number (note: required if enable_pr_comment is true)"
    default: 0

  "pr_comment_api_endpoint":
    description: "GitHub API endpoint URL (example: https://api.github.com or https://github.example.com/api/v3, note: required if enable_pr_comment is true)"
    default: "https://api.github.com"

  "pr_comment_no_send_success":
    description: "do not send a comment if no missing/dangling .meta found"
    default: false

`[1:]

	if expected != buf.String() {
		t.Log(buf.String())
		t.Error(cmp.Diff(expected, buf.String()))
	}
}

func TestWriteRunsAsYAML(t *testing.T) {
	buf := &bytes.Buffer{}
	runs := BuildMetadata().Runs

	if _, err := WriteRunsAsGHActionYAML(buf, runs); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := fmt.Sprintf(`
runs:
  using: "docker"
  image: "docker.pkg.github.com/dena/unity-meta-check/unity-meta-check-gh-action:%s"
  args:
    - "-inputs-json"
    - "{{ toJSON(github.inputs) }}"
    - "{{ github.workspace }}"
`[1:], version.Version)

	if expected != buf.String() {
		t.Log(buf.String())
		t.Error(cmp.Diff(expected, buf.String()))
	}
}
