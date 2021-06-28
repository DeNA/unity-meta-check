package inputs

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"testing"
)

func TestEventPayload(t *testing.T) {
	cases := map[string]struct {
		Path     string
		Expected EventPayload
	}{
		"pull request": {
			Path: "./testdata/pr-event-payload-example.json",
			Expected: EventPayload{
				PullRequest: &PullRequest{
					Number: 2,
				},
				Repository: &Repository{
					Name:  "Hello-World",
					Owner: User{Login: "Codertocat"},
				},
			},
		},
		"push": {
			Path: "./testdata/push-event-payload-example.json",
			Expected: EventPayload{
				Repository: &Repository{
					Name:  "Hello-World",
					Owner: User{Login: "Codertocat"},
				},
			},
		},
		"issue comment (comment to not pull request)": {
			Path: "./testdata/issue-comment-to-issue-payload-example.json",
			Expected: EventPayload{
				Issue:       &Issue{
					Number:      1,
				},
				Repository: &Repository{
					Name:  "Hello-World",
					Owner: User{Login: "Codertocat"},
				},
			},
		},
		"issue comment (comment to pull request)": {
			Path: "./testdata/issue-comment-to-pr-payload-example.json",
			Expected: EventPayload{
				Issue:       &Issue{
					Number:      1,
					PullRequest: &IssuePullRequest{URL: "https://api.github.com/repos/Codertocat/Hello-World/pulls/1"},
				},
				Repository: &Repository{
					Name:  "Hello-World",
					Owner: User{Login: "Codertocat"},
				},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			jsonBytes, err := os.ReadFile(c.Path)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			var payload EventPayload
			if err := json.Unmarshal(jsonBytes, &payload); err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if !reflect.DeepEqual(payload, c.Expected) {
				t.Error(cmp.Diff(c.Expected, payload))
			}
		})
	}
}
