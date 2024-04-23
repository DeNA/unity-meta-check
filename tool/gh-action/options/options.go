package options

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/pkg/errors"
)

type Options struct {
	Version bool
	Inputs  inputs.Inputs
	Env     inputs.ActionEnv
}

type Parser func(args []string, procInout cli.ProcessInout, env cli.Env) (*Options, error)

func NewParser() Parser {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) (*Options, error) {
		flags := flag.NewFlagSet("gh-action", flag.ContinueOnError)
		flags.SetOutput(procInout.Stderr)
		flags.Usage = func() {
			_, _ = fmt.Fprintln(flags.Output(), "usage: gh-action -inputs-json <json>")
			flags.PrintDefaults()
		}

		version := flags.Bool("version", false, "print version")
		inputsJSON := flags.String("inputs-json", "", `JSON string of "inputs" context value of GitHub Actions`)

		if err := flags.Parse(args); err != nil {
			return nil, err
		}

		if *version {
			return &Options{Version: *version}, nil
		}

		if flags.NArg() > 0 {
			return nil, fmt.Errorf("0 arguments required, but come %d arguments: %#v", flags.NArg(), flags.Args())
		}

		var unsafeInputs inputs.Inputs
		if err := json.Unmarshal([]byte(*inputsJSON), &unsafeInputs); err != nil {
			return nil, errors.Wrapf(err, "malformed JSON of inputs:\n%q", *inputsJSON)
		}

		return &Options{
			Inputs: unsafeInputs,
			Env:    inputs.GetActionEnv(env),
		}, nil
	}
}
