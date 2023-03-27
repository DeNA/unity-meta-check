package meta_test

import (
	"bytes"
	"testing"

	"github.com/DeNA/unity-meta-check/unity/meta"
	"github.com/google/go-cmp/cmp"
)

func TestDefaultoImporterGen_WriteTo(t *testing.T) {
	gen := meta.DefaultImporterGen{GUID: meta.ZeroGUID()}

	buf := &bytes.Buffer{}
	_, err := gen.WriteTo(buf)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := `fileFormatVersion: 2
guid: 00000000000000000000000000000000
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
