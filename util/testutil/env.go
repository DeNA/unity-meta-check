package testutil

import (
	"errors"
	"os"
	"strings"
)

type TestEnv struct {
	ApiEndpoint string
	Owner       string
	Repo        string
	Pull        string
	Token       string
}

type TestEnvError []string

func (t TestEnvError) Error() string {
	return strings.Join(t, "\n")
}

func GetTestEnv() (*TestEnv, error) {
	err := TestEnvError{}

	apiEndpoint := os.Getenv("UNITY_META_CHECK_GITHUB_API_ENDPOINT")
	if apiEndpoint == "" {
		err = append(err, "missing UNITY_META_CHECK_GITHUB_API_ENDPOINT")
	}

	owner := os.Getenv("UNITY_META_CHECK_GITHUB_OWNER")
	if owner == "" {
		err = append(err, "missing UNITY_META_CHECK_GITHUB_OWNER")
	}

	repo := os.Getenv("UNITY_META_CHECK_GITHUB_REPO")
	if repo == "" {
		err = append(err, "missing UNITY_META_CHECK_GITHUB_REPO")
	}

	pull := os.Getenv("UNITY_META_CHECK_GITHUB_PULL_NUMBER")
	if pull == "" {
		err = append(err, "missing UNITY_META_CHECK_GITHUB_PULL_NUMBER")
	}

	token := os.Getenv("UNITY_META_CHECK_GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("missing UNITY_META_CHECK_GITHUB_TOKEN")
	}

	return &TestEnv{
		apiEndpoint,
		owner,
		repo,
		pull,
		token,
	}, nil
}
