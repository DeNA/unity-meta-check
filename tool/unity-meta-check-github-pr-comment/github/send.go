package github

import (
	"bytes"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/markdown"
	"github.com/DeNA/unity-meta-check/unity/checker"
)

type Options struct {
	Tmpl          *l10n.Template
	SendIfSuccess bool
	Token         Token
	APIEndpoint   APIEndpoint
	Owner         Owner
	Repo          Repo
	PullNumber    PullNumber
}

type SendFunc func(result *checker.CheckResult, opts *Options) error

func NewSendFunc(postComment PullRequestCommentSender) SendFunc {
	return func(result *checker.CheckResult, opts *Options) error {
		buf := &bytes.Buffer{}
		if err := markdown.WriteMarkdown(result, opts.Tmpl, buf); err != nil {
			return err
		}

		if !result.Empty() || opts.SendIfSuccess {
			if err := postComment(opts.APIEndpoint, opts.Token, opts.Owner, opts.Repo, opts.PullNumber, buf.String()); err != nil {
				return err
			}
		}

		return nil
	}
}
