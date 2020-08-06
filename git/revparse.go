package git

import (
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/logging"
	"os/exec"
	"strings"
)

type RevParse func(repoDir string, options ...string) (string, error)

func NewRevParse(logger logging.Logger) RevParse {
	return func(repoDir string, options ...string) (string, error) {
		subcmdWithOpts := append([]string{"-c", "core.quotepath=false", "rev-parse"}, options...)
		logger.Debug(fmt.Sprintf("exec: git %s (on %q)", strings.Join(subcmdWithOpts, " "), repoDir))

		stdoutBuf := &bytes.Buffer{}
		stderrBuf := &bytes.Buffer{}

		cmd := exec.Command("git", subcmdWithOpts...)
		cmd.Dir = repoDir
		cmd.Stdout = stdoutBuf
		cmd.Stderr = stderrBuf

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("%s\nStderr:\n%s", err.Error(), stderrBuf.String())
		}

		return strings.TrimSpace(stdoutBuf.String()), nil
	}
}
