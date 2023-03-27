package meta

import (
	"fmt"
	"io"
	"strings"
)

type DefaultImporterGen struct {
	GUID *GUID
}

var _ Gen = MonoImporterGen{}

func (t DefaultImporterGen) WriteTo(writer io.Writer) (int64, error) {
	content := strings.TrimLeft(fmt.Sprintf(`
fileFormatVersion: 2
guid: %s
DefaultImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
`, t.GUID.String()), "\n")
	n, err := io.WriteString(writer, content)
	return int64(n), err
}
