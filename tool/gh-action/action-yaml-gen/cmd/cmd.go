package cmd

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/tool/gh-action/action-yaml-gen/yaml"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/pkg/errors"
	"io"
	"os"
)

func Main(args []string, procInout cli.ProcessInout, _ cli.Env) cli.ExitStatus {
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
	if err := writeYAML(actionYAMLPath); err != nil {
		_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
		return cli.ExitAbnormal
	}

	return cli.ExitNormal
}

func writeYAML(actionYAMLPath string) error {
	f, err := os.OpenFile(actionYAMLPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrapf(err, "cannot open: %q", actionYAMLPath)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if _, err := yaml.WriteTo(f); err != nil {
		return errors.Wrapf(err, "cannot write YAML to: %q", actionYAMLPath)
	}

	return nil
}
