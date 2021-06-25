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

type PullRequestCommentSender func(githubEndpoint APIEndpoint, token Token, owner Owner, repo Repo, pullNumber PullNumber, comment string) error

type PullRequestComment struct {
	Body string `json:"body"`
}

func NewPullRequestCommentSender(send HttpFunc, logger logging.Logger) PullRequestCommentSender {
	return func(githubEndpoint APIEndpoint, token Token, owner Owner, repo Repo, pullNumber PullNumber, comment string) error {
		apiUrl := &url.URL{
			Scheme: githubEndpoint.Scheme,
			Host:   githubEndpoint.Host,
			Path:   path.Join(githubEndpoint.Path, "repos", string(owner), string(repo), "issues", fmt.Sprintf("%d", pullNumber), "comments"),
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
		req.Header.Add("Accept", "application/vnd.github.v3+json")
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
