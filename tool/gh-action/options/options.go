package options

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/pkg/errors"
)

type Options struct {
	Version      bool
	UnsafeInputs inputs.Inputs
	Token        github.Token
	RootDirAbs   typedpath.RawPath
}

type Parser func(args []string, procInout cli.ProcessInout, env cli.Env) (*Options, error)

func NewParser(validateRootDirAbs options.RootDirAbsValidator) Parser {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) (*Options, error) {
		flags := flag.NewFlagSet("unity-meta-check-gh-action", flag.ContinueOnError)
		flags.SetOutput(procInout.Stderr)
		flags.Usage = func() {
			_, _ = fmt.Fprintln(flags.Output(), "usage: unity-meta-check-gh-action -inputs-json <json> <path>")
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

		if flags.NArg() != 1 {
			return nil, fmt.Errorf("must specify 1 argument as path to check, but come %d arguments: %#v", flags.NArg(), flags.Args())
		}

		unsafeRootDir := typedpath.NewRawPathUnsafe(flags.Arg(0))
		rootDirAbs, err := validateRootDirAbs(unsafeRootDir)
		if err != nil {
			return nil, err
		}

		var unsafeInputs inputs.Inputs
		if err := json.Unmarshal([]byte(*inputsJSON), &unsafeInputs); err != nil {
			return nil, errors.Wrapf(err, "malformed JSON of inputs:\n%q", *inputsJSON)
		}

		token, err := github.ValidateToken(env(options.GitHubTokenEnv))
		if err != nil {
			return nil, errors.Wrapf(err, "invalid environment variable: %s", options.GitHubTokenEnv)
		}

		return &Options{
			UnsafeInputs: unsafeInputs,
			Token:        token,
			RootDirAbs:   rootDirAbs,
		}, nil
	}
}
