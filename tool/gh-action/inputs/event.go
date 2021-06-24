package inputs

import (
	"encoding/json"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/pkg/errors"
	"os"
)

type ReadEventPayloadFunc func(path typedpath.RawPath) (*PullRequestEventPayload, error)

func ReadEventPayload(path typedpath.RawPath) (*PullRequestEventPayload, error) {
	payloadBytes, err := os.ReadFile(string(path))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read file: %q", path)
	}

	var payload PullRequestEventPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, errors.Wrapf(err, "cannot decode file: %q\n%s", path, string(payloadBytes))
	}

	return &payload, nil
}

// PullRequestEventPayload is a payload for pull request events.
// SEE: https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request
type PullRequestEventPayload struct {
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

// PullRequest is a payload for pull requests.
// SEE: https://docs.github.com/en/rest/reference/pulls#get-a-pull-request
type PullRequest struct {
	Number github.PullNumber `json:"number"`
}

// Repository is a payload for repository.
// SEE: https://docs.github.com/en/rest/reference/repos#get-a-repository
type Repository struct {
	URL   string      `json:"url"`
	Name  github.Repo `json:"name"`
	Owner User        `json:"owner"`
}

// User is a payload for users.
type User struct {
	Login github.Owner `json:"login"`
}
