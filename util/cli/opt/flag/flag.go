package flag

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/cli/opt"
)

func descSuffix(required bool, desc opt.Desc) string {
	if required {
		return fmt.Sprintf("%s (required)", desc)
	}
	return fmt.Sprintf("%s (optional)", desc)
}

func DefineString(f *flag.FlagSet, o opt.StringOption) *string {
	return f.String(string(o.Name()), o.DefaultValue, descSuffix(o.Required(), o.Desc()))
}

func DefineBool(f *flag.FlagSet, o opt.BoolOption) *bool {
	return f.Bool(string(o.Name()), o.DefaultValue, descSuffix(o.Required(), o.Desc()))
}
