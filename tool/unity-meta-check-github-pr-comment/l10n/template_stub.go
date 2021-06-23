package l10n

import "github.com/DeNA/unity-meta-check/util/typedpath"

func StubTemplateFileReader(tmpl *Template, err error) TemplateFileReader {
	return func(_ typedpath.RawPath) (*Template, error) {
		return tmpl, err
	}
}
