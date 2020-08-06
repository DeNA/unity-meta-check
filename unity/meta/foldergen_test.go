package meta

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFolderGen_WriteTo(t *testing.T) {
	meta := DefaultImporterFolderGen{ZeroGUID()}

	buf := &bytes.Buffer{}
	_, err := meta.WriteTo(buf)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := `fileFormatVersion: 2
guid: 00000000000000000000000000000000
folderAsset: yes
DefaultImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
`
	if actual != expected {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
