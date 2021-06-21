package markdown

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/cli/opt"
	"io"
)

func WriteTo(w io.Writer, opts ...opt.Option) (i int64, err error) {
	opt.Sort(opts)

	j, err := fmt.Fprintln(w, `
| Option | Description | Required or Default Value |
|:-------|:------------|:--------------------------|`[1:])
	i += int64(j)
	if err != nil {
		return
	}

	for _, o := range opts {
		var required string
		if o.Required() {
			required = "required"
		} else {
			required = fmt.Sprintf("optional (default: `%s`)", o.DefaultValueText())
		}

		var oText string
		oText, err = optionText(o)
		if err != nil {
			return
		}

		var j int
		j, err = fmt.Fprintf(w, "| `%s` | %s | %s |\n", oText, o.Desc(), required)
		i += int64(j)
		if err != nil {
			return
		}
	}

	return
}

func optionText(o opt.Option) (string, error) {
	switch o.(type) {
	case opt.StringOption:
		return fmt.Sprintf("--%s <string>", o.Name()), nil
	case opt.BoolOption:
		return fmt.Sprintf("--%s", o.Name()), nil
	default:
		return "", fmt.Errorf("unknown option type: %#v", o)
	}
}