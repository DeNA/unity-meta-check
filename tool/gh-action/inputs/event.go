package inputs

import (
	"encoding/json"
	"fmt"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/pkg/errors"
	"os"
)

type ReadEventPayloadFunc func(path typedpath.RawPath) (*EventPayload, error)

func NewReadEventPayload(logger logging.Logger) ReadEventPayloadFunc {
	return func(path typedpath.RawPath) (*EventPayload, error) {
		payloadBytes, err := os.ReadFile(string(path))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read file: %q", path)
		}

		logger.Debug(fmt.Sprintf("event=%s", string(payloadBytes)))

		var payload EventPayload
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrapf(err, "cannot decode file: %q\n%s", path, string(payloadBytes))
		}

		return &payload, nil
	}
}

type EventPayload struct {
	// PullRequest is a payload for pull request related events.
	// SEE: https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request
	PullRequest *PullRequest `json:"pull_request,omitempty"`

	// Issue is a payload for issue related events.
	// SEE: https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#issue_comment
	Issue *Issue `json:"issue"`

	// Repository is the repository where the event occurred.
	Repository *Repository `json:"repository,omitempty"`
}

// PullRequest is a payload for GitHub pull requests.
// SEE: https://docs.github.com/en/rest/reference/pulls#get-a-pull-request
type PullRequest struct {
	Number github.PullNumber `json:"number"`
}

// Repository is a payload for GitHub repositories.
// SEE: https://docs.github.com/en/rest/reference/repos#get-a-repository
type Repository struct {
	Name  github.Repo `json:"name"`
	Owner User        `json:"owner"`
}

type IssueNumber int

// Issue is a payload for GitHub issue related events.
// SEE: https://docs.github.com/en/rest/reference/issues#get-an-issue
type Issue struct {
	Number IssueNumber `json:"number"`

	// PullRequest exists if the issue is a pull request
	PullRequest *IssuePullRequest `json:"pull_request,omitempty"`
}

func (i Issue) GetPullNumber() (github.PullNumber, error) {
	if i.PullRequest != nil {
		return github.PullNumber(i.Number), nil
	}
	return 0, fmt.Errorf("not a pull request: %d", i.Number)
}

type IssuePullRequest struct {
	URL string `json:"url"`
}

// User is a payload for users.
type User struct {
	Login github.Owner `json:"login"`
}
