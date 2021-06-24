package inputs

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"testing"
)

func TestEventPayload(t *testing.T) {
	jsonBytes, err := os.ReadFile("./testdata/webhook-payload-example.json")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	var payload PullRequestEventPayload
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := PullRequestEventPayload{
		PullRequest: PullRequest{
			Number: 2,
		},
		Repository: Repository{
			URL:   "https://api.github.com/repos/Codertocat/Hello-World",
			Name:  "Hello-World",
			Owner: User{Login: "Codertocat"},
		},
	}
	if !reflect.DeepEqual(payload, expected) {
		t.Error(cmp.Diff(expected, payload))
	}
}
