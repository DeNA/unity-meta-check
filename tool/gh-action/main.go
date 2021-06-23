package main

import (
	"github.com/DeNA/unity-meta-check/tool/gh-action/cmd"
	"github.com/DeNA/unity-meta-check/util/cli"
	"os"
)

func main() {
	exitStatus := cmd.Main(os.Args[1:], cli.GetProcessInout(), cli.NewEnv())
	os.Exit(int(exitStatus))
}
