package yaml

import (
	"github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/version"
)

func BuildMetadata() GHActionsMetadata {
	return GHActionsMetadata{
		Name:        "Unity meta files check",
		Author:      "DeNA Co., Ltd.",
		Description: "Check missing/dangling meta files",
		Inputs: []InputDef{
			StringInputDef{
				name:         "log_level",
				desc:         "log level for unity-meta-check toolchains (available: ERROR/WARN/INFO/DEBUG)",
				defaultValue: "INFO",
			},
			StringInputDef{
				name:         "target_type",
				desc:         "target type for unity-meta-check (available: auto-detect/unity-project/unity-project-sub-dir/upm-package)",
				defaultValue: "auto-detect",
			},
			BoolInputDef{
				name:         "ignore_dangling",
				desc:         "ignore dangling .meta",
				defaultValue: false,
			},
			BoolInputDef{
				name:         "ignore_case",
				desc:         "do not treat case of file paths",
				defaultValue: false,
			},
			BoolInputDef{
				name:         "ignore_submodules_and_nested",
				desc:         "ignore git submodules and nesting repositories (note: this is RECOMMENDED but not enabled by default because it can cause to miss problems in submodules or nesting repositories)",
				defaultValue: false,
			},
			StringInputDef{
				name:         "ignored_file_path",
				desc:         "path to .meta-check-ignore",
				defaultValue: string(options.IgnoreFileBasename),
			},
			BoolInputDef{
				name:         "enable_autofix",
				desc:         "enable autofix (note: it can repair very limited file types)",
				defaultValue: false,
			},
			BoolInputDef{
				name:         "autofix_missing",
				desc:         "fix missing .meta",
				defaultValue: true,
			},
			BoolInputDef{
				name:         "autofix_dangling",
				desc:         "fix dangling .meta",
				defaultValue: true,
			},
			StringInputDef{
				name:         "autofix_globs",
				desc:         "glob match directories to where to allow do autofix (note: required if enable_autofix is true)",
				defaultValue: "",
			},
			BoolInputDef{
				name:         "enable_junit",
				desc:         "enable JUnit-style reporting",
				defaultValue: false,
			},
			StringInputDef{
				name:         "junit_xml_path",
				desc:         "file path for generated JUnit test reports (note: required if enable_junit is true)",
				defaultValue: "",
			},
			BoolInputDef{
				name:         "enable_pr_comment",
				desc:         "enable reporting via GitHub Pull Request Comments",
				defaultValue: false,
			},
			BoolInputDef{
				name:         "pr_comment_lang",
				desc:         "language code for GitHub Pull Request Comments (available: en/ja, note: cannot specify both lang and pr_comment_tmpl_file)",
				defaultValue: false,
			},
			StringInputDef{
				name:         "pr_comment_tmpl_file",
				desc:         "file path to custom template file for GitHub Pull Request Comments (note: cannot specify both lang and pr_comment_tmpl_file)",
				defaultValue: "",
			},
			StringInputDef{
				name:         "pr_comment_owner",
				desc:         "owner of the GitHub repository (note: required if enable_pr_comment is true)",
				defaultValue: "",
			},
			StringInputDef{
				name:         "pr_comment_repo",
				desc:         "name of the GitHub repository (note: required if enable_pr_comment is true)",
				defaultValue: "",
			},
			IntInputDef{
				name:         "pr_comment_pull",
				desc:         "pull request number (note: required if enable_pr_comment is true)",
				defaultValue: 0,
			},
			StringInputDef{
				name:         "pr_comment_api_endpoint",
				desc:         "GitHub API endpoint URL (example: https://api.github.com or https://github.example.com/api/v3, note: required if enable_pr_comment is true)",
				defaultValue: "https://api.github.com",
			},
			BoolInputDef{
				name:         "pr_comment_no_send_success",
				desc:         "do not send a comment if no missing/dangling .meta found",
				defaultValue: false,
			},
		},
		Runs: DockerAction{
			Image: ImageWithTag(version.Version),
			Args: []string{
				`-inputs-json`,
				`{{ toJSON(github.inputs) }}`,
				`{{ github.workspace }}`,
			},
		},
		BrandingIcon:  BrandingIconCheck,
		BrandingColor: BrandingColorBlue,
	}
}

const (
	BrandingIconCheck BrandingIcon  = "check"
	BrandingColorBlue BrandingColor = "blue"
)

type GHActionsMetadata struct {
	// Name is the name of your action. GitHub displays the name in the Actions tab to help visually identify actions in each job.
	// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#name
	Name string

	// Author is the name of the action's author.
	// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#author
	Author string

	// Description is a short description of the action.
	// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#description
	Description string

	// Inputs is Input parameters allow you to specify data that the action expects to use during runtime.
	// GitHub stores input parameters as environment variables. Input ids with uppercase letters are converted to lowercase during runtime.
	// We recommended using lowercase input ids.
	// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#inputs
	Inputs []InputDef

	// Runs configures the image used for the Docker action.
	// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#runs-for-docker-actions
	Runs DockerAction

	// BrandingIcon is the name of the Feather icon to use.
	// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#brandingicon
	BrandingIcon BrandingIcon

	// BrandingColor is the background color of the badge. Can be one of: white, yellow, blue, green, orange, red, purple, or gray-dark.
	// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#brandingcolor
	BrandingColor BrandingColor
}
