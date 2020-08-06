package github

import (
	"github.com/DeNA/unity-meta-check/util/logging"
	"net/url"
	"os"
	"testing"
)

func TestNewPullRequestCommentSender(t *testing.T) {
	spyLogger := logging.SpyLogger()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skipf("need GITHUB_TOKEN")
		return
	}

	endpoint, err := url.Parse("https://api.github.com/")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	send := NewPullRequestCommentSender(endpoint, token, NewHttp(), spyLogger)
	if err := send("dena", "unity-meta-check-playground", 1, "TEST FROM unity-meta-check-github-pr-comment"); err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}
}
