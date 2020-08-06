package meta

import (
	"fmt"
	"io"
	"strings"
)

type TextScriptImporterGen struct {
	GUID *GUID
}

var _ Gen = TextScriptImporterGen{}

func (t TextScriptImporterGen) WriteTo(writer io.Writer) (int64, error) {
	content := strings.TrimLeft(fmt.Sprintf(`
fileFormatVersion: 2
guid: %s
TextScriptImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
`, t.GUID.String()), "\n")
	n, err := io.WriteString(writer, content)
	return int64(n), err
}

