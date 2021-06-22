package cmd

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/tool/gh-action/action-yaml-gen/yaml"
	"github.com/DeNA/unity-meta-check/util/cli"
	"io"
	"os"
)

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		if len(args) != 1 {
			_, _ = fmt.Fprintf(procInout.Stderr, "error: must specify 1 argument, but come %d arguments\n", len(args))
			return cli.ExitAbnormal
		}

		flags := flag.NewFlagSet("gh-action-yaml-gen", flag.ContinueOnError)
		flags.SetOutput(io.Discard)

		if err := flags.Parse(args); err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, "usage: gh-action-yaml-gen <path>")
			return cli.ExitAbnormal
		}

		actionYAMLPath := args[0]

		f, err := os.OpenFile(actionYAMLPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			_, _ = fmt.Fprintf(procInout.Stderr, "cannot open path for action.yml: %q\n%s", actionYAMLPath, err.Error())
			return cli.ExitAbnormal
		}
		defer f.Close()

		if _, err := yaml.WriteTo(f); err != nil {
			_, _ = fmt.Fprintf(procInout.Stderr, "cannot write to action.yml: %q\n%s", actionYAMLPath, err.Error())
			return cli.ExitAbnormal
		}

		return cli.ExitNormal
	}
}
