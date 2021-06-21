package ghactions

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/cli/opt"
	"io"
)

func WriteTo(w io.Writer, flags ...opt.Option) (i int64, err error) {
	opt.Sort(flags)

	for _, f := range flags {
		var j int
		if f.Required() {
			j, err = fmt.Fprintf(w, `
  %q:
    description: %q
    required: true
`[1:], f.Name(), f.Desc())
		} else {
			var defVal string
			defVal, err = defaultValue(f)
			if err != nil {
				return
			}

			j, err = fmt.Fprintf(w, `
  %q:
    description: %q
    required: false
    default: %s
`[1:], f.Name(), f.Desc(), defVal)
		}

		i += int64(j)
		if err != nil {
			return
		}
	}

	return
}

func defaultValue(f opt.Option) (string, error) {
	switch f.(type) {
	case opt.StringOption:
		return fmt.Sprintf("%q", f.(opt.StringOption).DefaultValue), nil
	case opt.BoolOption:
		return fmt.Sprintf("%v", f.(opt.BoolOption).DefaultValue), nil
	default:
		return "", fmt.Errorf("unknown opt.Option type: %#v", f)
	}
}
