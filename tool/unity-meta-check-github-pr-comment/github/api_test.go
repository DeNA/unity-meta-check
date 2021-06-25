package github

import (
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"net/url"
	"testing"
)

func TestNewPullRequestCommentSender(t *testing.T) {
	spyLogger := logging.SpyLogger()

	testEnv, err := testutil.GetTestEnv()
	if err != nil {
		t.Log(err.Error())
		t.Skip("no environment variables for tests")
		return
	}
	token := testEnv.Token

	endpoint, err := url.Parse("https://api.github.com")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	send := NewPullRequestCommentSender(NewHttp(), spyLogger)
	if err := send(endpoint, Token(token), "dena", "unity-meta-check-playground", 1, "TEST FROM unity-meta-check-github-pr-comment"); err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}
}
