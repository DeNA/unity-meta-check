package cli

import (
	"io"
	"os"
)

type ProcessInout struct {
	Stdin  io.Reader
	Stdout io.WriteCloser
	Stderr io.WriteCloser
}

func GetProcessInout() ProcessInout {
	return ProcessInout{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

