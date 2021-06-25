package junit

import (
	"errors"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"time"
)

func StubWriteToFileFunc(err error) WriteToFileFunc {
	return func(_ *checker.CheckResult, _ time.Time, _ typedpath.RawPath) error {
		return err
	}
}

type WriteToFileCallArgs struct {
	CheckResult *checker.CheckResult
	StartTime   time.Time
	OutPath     typedpath.RawPath
}

func SpyWriteToFileFunc(inherited WriteToFileFunc, callArgs *[]WriteToFileCallArgs) WriteToFileFunc {
	if inherited == nil {
		inherited = StubWriteToFileFunc(errors.New("SPY_WRITE_TO_FILE_FUNC"))
	}
	return func(result *checker.CheckResult, startTime time.Time, outPath typedpath.RawPath) error {
		*callArgs = append(*callArgs, WriteToFileCallArgs{
			CheckResult: result,
			StartTime:   startTime,
			OutPath:     outPath,
		})
		return inherited(result, startTime, outPath)
	}
}