package meta

import (
	"fmt"
	"io"
	"strings"
)

type DefaultImporterFolderGen struct {
	GUID *GUID
}

var _ Gen = DefaultImporterFolderGen{}

func (f DefaultImporterFolderGen) WriteTo(writer io.Writer) (int64, error) {
	content := strings.TrimLeft(fmt.Sprintf(`
fileFormatVersion: 2
guid: %s
folderAsset: yes
DefaultImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
`, f.GUID.String()), "\n")
	n, err := io.WriteString(writer, content)
	return int64(n), err
}
