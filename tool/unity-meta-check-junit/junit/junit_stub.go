package junit

import (
	"errors"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubWriteToFileFunc(err error) WriteToFileFunc {
	return func(_ *checker.CheckResult, _ typedpath.RawPath) error {
		return err
	}
}

type WriteToFileCallArgs struct {
	CheckResult *checker.CheckResult
	OutPath     typedpath.RawPath
}

func SpyWriteToFileFunc(inherited WriteToFileFunc, callArgs *[]WriteToFileCallArgs) WriteToFileFunc {
	if inherited == nil {
		inherited = StubWriteToFileFunc(errors.New("SPY_WRITE_TO_FILE_FUNC"))
	}
	return func(result *checker.CheckResult, outPath typedpath.RawPath) error {
		*callArgs = append(*callArgs, WriteToFileCallArgs{
			CheckResult: result,
			OutPath:     outPath,
		})
		return inherited(result, outPath)
	}
}
