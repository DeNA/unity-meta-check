package meta

import (
	"fmt"
	"io"
	"strings"
)

type MonoImporterGen struct {
	GUID *GUID
}

var _ Gen = MonoImporterGen{}

func (t MonoImporterGen) WriteTo(writer io.Writer) (int64, error) {
	content := strings.TrimLeft(fmt.Sprintf(`
fileFormatVersion: 2
guid: %s
MonoImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
  serializedVersion: 2
  defaultReferences: []
  executionOrder: 0
  icon: {instanceID: 0}
`, t.GUID.String()), "\n")
	n, err := io.WriteString(writer, content)
	return int64(n), err
}
