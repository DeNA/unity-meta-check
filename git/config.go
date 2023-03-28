package git

import (
	"bufio"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/logging"
	"io"
	"os/exec"
	"sync"
)

type GlobalConfig func(stdoutWriter io.WriteCloser, args ...string) error

func NewGlobalConfig(logger logging.Logger) GlobalConfig {
	return func(stdoutWriter io.WriteCloser, args ...string) error {
		configWithOpts := append([]string{"config", "--global"}, args...)
		cmd := exec.Command("git", configWithOpts...)
		cmd.Stdout = stdoutWriter
		defer func() { _ = stdoutWriter.Close() }()

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
