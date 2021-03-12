package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func NewCommandWithCPUProfile(cmd Command) Command {
	profFile, err := newProfFile()
	if err != nil {
		panic(err.Error())
	}

	return func(args []string, procInout ProcessInout, env Env) ExitStatus {
		if err := pprof.StartCPUProfile(profFile); err != nil {
			panic(err.Error())
		}

		exitStatus := cmd(args, procInout, env)

		pprof.StopCPUProfile()
		_ = profFile.Close()
		return exitStatus
	}
}

func NewCommandWithHeapProfile(cmd Command) Command {
	profFile, err := newProfFile()
	if err != nil {
		panic(err.Error())
	}

	return func(args []string, procInout ProcessInout, env Env) ExitStatus {
		exitStatus := cmd(args, procInout, env)

		if err := pprof.Lookup("heap").WriteTo(profFile, 0); err != nil {
			panic(err.Error())
		}
		_ = profFile.Close()
		return exitStatus
	}
}

func newProfFile() (io.WriteCloser, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}

	profPath := filepath.Join(tmpDir, "unity-meta-check.prof")
	fmt.Printf("profile path: %s\n", profPath)

	profFile, err := os.OpenFile(profPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}
	return profFile, nil
}
