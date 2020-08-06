package git

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io"
	"strings"
)

func FakeLsFiles(lsFiles func(repoDir typedpath.RawPath)([]string, error)) LsFiles {
	return func(repoDir typedpath.RawPath, options []string, stdoutWriter io.WriteCloser) error {
		files, err := lsFiles(repoDir)
		if err != nil {
			return err
		}
		_, _ = io.WriteString(stdoutWriter, strings.Join(files, "\n"))
		_ = stdoutWriter.Close()
		return err
	}
}

func StubLsFiles(files []string, err error) LsFiles {
	return StubRawLsFiles(strings.Join(files, "\n") + "\n", err)
}

func StubRawLsFiles(stdout string, err error) LsFiles {
	return func(repoDir typedpath.RawPath, options []string, stdoutWriter io.WriteCloser) error {
		_, _ = io.WriteString(stdoutWriter, stdout)
		_ = stdoutWriter.Close()
		return err
	}
}
