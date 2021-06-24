package typedpath

import "testing"

func TestNewRawPath(t *testing.T) {
	cases := map[string]struct {
		BaseNames []BaseName
		Expected  RawPath
	}{
		"empty": {
			BaseNames: nil,
			Expected: NewSlashPathUnsafe("/").ToRaw(),
		},
		"several paths": {
			BaseNames: []BaseName{"path", "to", "file"},
			Expected: NewSlashPathUnsafe("/path/to/file").ToRaw(),
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			actual := NewRootRawPath(c.BaseNames...)

			if actual != c.Expected {
				t.Errorf("want %q, got %q", c.Expected, actual)
			}
		})
	}
}
