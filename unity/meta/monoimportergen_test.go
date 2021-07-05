package meta_test

import (
	"bytes"
	"github.com/DeNA/unity-meta-check/unity/meta"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestMonoImporterGen_WriteTo(t *testing.T) {
	gen := meta.MonoImporterGen{GUID: meta.ZeroGUID()}

	buf := &bytes.Buffer{}
	_, err := gen.WriteTo(buf)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := `fileFormatVersion: 2
guid: 00000000000000000000000000000000
MonoImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
  serializedVersion: 2
  defaultReferences: []
  executionOrder: 0
  icon: {instanceID: 0}
`
	if actual != expected {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
