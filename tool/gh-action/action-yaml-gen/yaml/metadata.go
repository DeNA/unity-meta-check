package yaml

import (
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
				desc:         "Log level for unity-meta-check toolchains (available: ERROR/WARN/INFO/DEBUG)",
				defaultValue: "INFO",
			},
			StringInputDef{
				name:         "target_path",
				desc:         "Target directory path to check. Relative path is treated as a path from $GITHUB_WORKSPACE (default: .)",
				defaultValue: ".",
			},
			StringInputDef{
				name:         "target_type",
				desc:         "Target type for unity-meta-check (available: auto-detect/unity-project/unity-project-sub-dir/upm-package)",
				defaultValue: "auto-detect",
			},
			BoolInputDef{
				name:         "ignore_dangling",
				desc:         "Ignore dangling .meta",
				defaultValue: false,
			},
			BoolInputDef{
				name:         "ignore_case",
				desc:         "Do not treat case of file paths",
				defaultValue: false,
			},
			BoolInputDef{
				name:         "ignore_submodules_and_nested",
				desc:         "Ignore git submodules and nesting repositories (note: this is RECOMMENDED but not enabled by default because it can cause to miss problems in submodules or nesting repositories)",
				defaultValue: false,
			},
			StringInputDef{
				name:         "ignored_file_path",
				desc:         "Path to .meta-check-ignore (note: keep empty if .meta-check-ignore does not exist)",
				defaultValue: "",
			},
			BoolInputDef{
				name:         "enable_autofix",
				desc:         "Enable autofix (note: it can repair very limited file types)",
				defaultValue: false,
			},
			StringInputDef{
				name:         "autofix_globs",
				desc:         "Comma separated globs that match directories to where to allow autofix (example: use '*' to allow autofix all files, note: required if enable_autofix is true)",
				defaultValue: "*",
			},
			BoolInputDef{
				name:         "enable_junit",
				desc:         "Enable JUnit-style reporting",
				defaultValue: false,
			},
			StringInputDef{
				name:         "junit_xml_path",
				desc:         "File path for generated JUnit test reports (note: required if enable_junit is true)",
				defaultValue: "",
			},
			BoolInputDef{
				name:         "enable_pr_comment",
				desc:         "Enable reporting via GitHub Pull Request Comments",
				defaultValue: false,
			},
			StringInputDef{
				name:         "pr_comment_lang",
				desc:         "Language code for GitHub Pull Request Comments (available: en/ja, note: cannot specify both lang and pr_comment_tmpl_file)",
				defaultValue: "en",
			},
			StringInputDef{
				name:         "pr_comment_tmpl_file",
				desc:         "File path to custom template file for GitHub Pull Request Comments (note: cannot specify both lang and pr_comment_tmpl_file)",
				defaultValue: "",
			},
			BoolInputDef{
				name:         "pr_comment_send_success",
				desc:         "Send a comment if no missing/dangling .meta found",
				defaultValue: false,
			},
		},
		Runs: DockerAction{
			Image: ImageWithTag(version.Version),
			Args: []string{
				`-inputs-json`,
				`${{ toJSON(inputs) }}`,
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
