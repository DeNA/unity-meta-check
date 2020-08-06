package unity

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestReadManifestJson(t *testing.T) {
	json := []byte(`{
  "scopedRegistries": [],
  "dependencies": {
    "foo": "1.2.3",
    "bar.baz": "file:../Bar/Buz"
  }
}`)

	actual, err := parseManifestJson(json)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := &ManifestJson{Dependencies: map[string]string{
		"foo": "1.2.3",
		"bar.baz": "file:../Bar/Buz",
	}}

	if !reflect.DeepEqual(actual, expected) {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
