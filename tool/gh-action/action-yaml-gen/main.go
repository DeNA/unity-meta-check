package main

import (
	cmd2 "github.com/DeNA/unity-meta-check/tool/gh-action/action-yaml-gen/cmd"
	"github.com/DeNA/unity-meta-check/util/cli"
	"os"
)

func main() {
	main := cmd2.NewMain()
	exitStatus := main(os.Args[1:], cli.GetProcessInout(), cli.NewEnv())
	os.Exit(int(exitStatus))
}
