package github

import (
	"errors"
	"fmt"
	"net/url"
)

type APIEndpoint *url.URL
type Token string
type Owner string
type Repo string
type PullNumber int

func ValidateOwner(unsafeOwner string) (Owner, error) {
	if unsafeOwner == "" {
		return "", errors.New("owner must be not empty")
	}
	return Owner(unsafeOwner), nil
}

func ValidateRepo(unsafeRepo string) (Repo, error) {
	if unsafeRepo == "" {
		return "", errors.New("repo must be not empty")
	}
	return Repo(unsafeRepo), nil
}

func ValidatePullNumber(unsafePullNumber int) (PullNumber, error) {
	if unsafePullNumber <= 0 {
		return 0, fmt.Errorf("pull number must be a positive integer: %d", unsafePullNumber)
	}
	return PullNumber(unsafePullNumber), nil
}

func ValidateToken(unsafeToken string) (Token, error) {
	if unsafeToken == "" {
		return "", fmt.Errorf("GitHub Personal Token must not be empty")
	}
	return Token(unsafeToken), nil
}

func ValidateAPIEndpoint(unsafeAPIEndpoint string) (APIEndpoint, error) {
	apiEndpoint, err := url.Parse(unsafeAPIEndpoint)
	if err != nil {
		return nil, err
	}
	return apiEndpoint, nil
}