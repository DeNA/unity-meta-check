package main

import (
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/cmd"
	"github.com/DeNA/unity-meta-check/util/cli"
	"os"
)

func main() {
	main := cmd.NewMain()
	exitStatus := main(os.Args[1:], cli.GetProcessInout(), cli.NewEnv())
	os.Exit(int(exitStatus))
}
