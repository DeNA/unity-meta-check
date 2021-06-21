package flag

import (
	"bytes"
	"flag"
	"github.com/DeNA/unity-meta-check/util/cli/opt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDefineHelp(t *testing.T) {
	t.Run("required string option (easies case)", func(t *testing.T) {
		buf := &bytes.Buffer{}
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.SetOutput(buf)

		o := opt.NewRequiredStringOption("req-str", "required string option")
		DefineString(flags, o)

		flags.PrintDefaults()

		expected := `
  -req-str string
    	required string option (required)
`[1:]

		if expected != buf.String() {
			t.Error(cmp.Diff(expected, buf.String()))
		}
	})

	t.Run("optional string option (easies case)", func(t *testing.T) {
		buf := &bytes.Buffer{}
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.SetOutput(buf)

		o := opt.NewOptionalStringOption("opt-str", "optional string option", "DEFAULT")
		DefineString(flags, o)

		flags.PrintDefaults()

		expected := `
  -opt-str string
    	optional string option (optional) (default "DEFAULT")
`[1:]

		if expected != buf.String() {
			t.Error(cmp.Diff(expected, buf.String()))
		}
	})

	t.Run("required bool option (easies case)", func(t *testing.T) {
		buf := &bytes.Buffer{}
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.SetOutput(buf)

		o := opt.NewRequiredBoolOption("req-bool", "required bool option")
		DefineBool(flags, o)

		flags.PrintDefaults()

		expected := `
  -req-bool
    	required bool option (required)
`[1:]

		if expected != buf.String() {
			t.Error(cmp.Diff(expected, buf.String()))
		}
	})

	t.Run("optional bool option (easies case)", func(t *testing.T) {
		buf := &bytes.Buffer{}
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.SetOutput(buf)

		o := opt.NewOptionalBoolOption("opt-bool", "optional bool option", true)
		DefineBool(flags, o)

		flags.PrintDefaults()

		expected := `
  -opt-bool
    	optional bool option (optional) (default true)
`[1:]

		if expected != buf.String() {
			t.Error(cmp.Diff(expected, buf.String()))
		}
	})
}
