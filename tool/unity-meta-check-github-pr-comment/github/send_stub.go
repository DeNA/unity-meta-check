package github

import (
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/pkg/errors"
)

func StubSendFunc(err error) SendFunc {
	return func(_ *checker.CheckResult, _ *Options) error {
		return err
	}
}

type SendFuncCallArgs struct {
	Result *checker.CheckResult
	Options *Options
}

func SpySendFunc(inherited SendFunc, callArgs *[]SendFuncCallArgs) SendFunc {
	if inherited == nil {
		inherited = StubSendFunc(errors.New("SPY_SEND_FUNC"))
	}
	return func(result *checker.CheckResult, opts *Options) error {
		*callArgs = append(*callArgs, SendFuncCallArgs{
			Result:  result,
			Options: opts,
		})
		return inherited(result, opts)
	}
}
