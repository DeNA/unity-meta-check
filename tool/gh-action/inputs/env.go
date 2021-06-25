package inputs

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"strings"
)

type ActionEnv struct {
	GitHubToken github.Token
	EventPath   typedpath.RawPath
	Workspace   typedpath.RawPath
	APIURL      string
}

// SEE: https://docs.github.com/en/actions/reference/environment-variables#default-environment-variables
const (
	// GitHubEventPath is the path of the file with the complete webhook event payload. For example, /github/workflow/event.json.
	GitHubEventPath = "GITHUB_EVENT_PATH"

	// GitHubWorkspace is the GitHub workspace directory path. The workspace directory is a copy of your repository
	// if your workflow uses the actions/checkout action. If you don't use the actions/checkout action,
	// the directory will be empty. For example, /home/runner/work/my-repo-name/my-repo-name.
	GitHubWorkspace = "GITHUB_WORKSPACE"

	// GitHubAPIURL returns the API URL. For example: https://api.github.com.
	GitHubAPIURL = "GITHUB_API_URL"
)

func GetActionEnv(env cli.Env) ActionEnv {
	gitHubToken := env(options.GitHubTokenEnv)
	eventPath := env(GitHubEventPath)
	workspace := env(GitHubWorkspace)
	apiURL := env(GitHubAPIURL)

	return ActionEnv{
		GitHubToken: github.Token(gitHubToken),
		EventPath:   typedpath.NewRawPathUnsafe(eventPath),
		Workspace:   typedpath.NewRawPathUnsafe(workspace),
		APIURL:      apiURL,
	}
}

func MaskedActionEnv(env ActionEnv) string {
	return fmt.Sprintf(`
TOKEN=%q (len=%d)
GITHUB_EVENT_PATH=%q
GITHUB_WORKSPACE=%q
GITHUB_API_URL=%q`[1:], strings.Repeat("*", len(env.GitHubToken)), len(env.GitHubToken), env.EventPath, env.Workspace, env.APIURL)
}
