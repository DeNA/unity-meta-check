package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/logging"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type PullRequestCommentSender func(owner string, repo string, prNumber uint, comment string) error

type PullRequestComment struct {
	Body string `json:"body"`
}

func NewPullRequestCommentSender(githubEndpoint *url.URL, token string, send HttpFunc, logger logging.Logger) PullRequestCommentSender {
	return func(owner string, repo string, prNumber uint, comment string) error {
		apiUrl := &url.URL{
			Scheme: githubEndpoint.Scheme,
			Host:   githubEndpoint.Host,
			Path:   path.Join(githubEndpoint.Path, "repos", owner, repo, "issues", fmt.Sprintf("%d", prNumber), "comments"),
		}
		logger.Debug(fmt.Sprintf("url: %s", apiUrl.String()))

		body, err := json.Marshal(PullRequestComment{Body: comment})
		if err != nil {
			return err
		}
		logger.Debug(fmt.Sprintf("comment body: %s", body))

		req, err := http.NewRequest(http.MethodPost, apiUrl.String(), bytes.NewReader(body))
		if err != nil {
			return err
		}

		logger.Debug(fmt.Sprintf("token: %s", strings.Repeat("*", len(token))))
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
		req.Header.Add("Content-Type", "application/json")

		res, err := send(req)
		if err != nil {
			return err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusCreated {
			return fmt.Errorf("unexpected status code: %d\n%s", res.StatusCode, string(resBody))
		}

		logger.Debug(fmt.Sprintf("response: %s", resBody))
		return nil
	}
}
