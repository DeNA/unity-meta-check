package git

import (
	"bufio"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io"
	"os/exec"
	"strings"
	"sync"
)

type LsFiles func(repoDir typedpath.RawPath, options []string, stdoutWriter io.WriteCloser) error

func NewLsFiles(logger logging.Logger) LsFiles {
	return func(repoDir typedpath.RawPath, options []string, stdoutWriter io.WriteCloser) error {
		subcmdWithOpts := append([]string{"-c", "core.quotepath=false", "ls-files"}, options...)
		logger.Debug(fmt.Sprintf("exec: git %s (on %q)", strings.Join(subcmdWithOpts, " "), repoDir))

		cmd := exec.Command("git", subcmdWithOpts...)
		cmd.Dir = string(repoDir)
		cmd.Stdout = stdoutWriter
		defer func(){ _ = stdoutWriter.Close() }()

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				logger.Debug(fmt.Sprintf("stderr: %s", scanner.Text()))
			}
		}()

		err = cmd.Run()
		wg.Wait()
		return err
	}
}
