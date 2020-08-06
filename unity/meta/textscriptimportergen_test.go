package meta

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTextScriptImporterGen_WriteTo(t *testing.T) {
	meta := TextScriptImporterGen{ZeroGUID()}

	buf := &bytes.Buffer{}
	_, err := meta.WriteTo(buf)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := `fileFormatVersion: 2
guid: 00000000000000000000000000000000
TextScriptImporter:
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
