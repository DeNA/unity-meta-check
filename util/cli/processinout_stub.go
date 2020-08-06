package cli

import (
	"github.com/DeNA/unity-meta-check/util/testutil"
)

func AnyProcInout() ProcessInout {
	return ProcessInout{
		Stdin:  &testutil.ErrorReadCloserStub{},
		Stdout: &testutil.NullWriteCloser{},
		Stderr: &testutil.NullWriteCloser{},
	}
}
