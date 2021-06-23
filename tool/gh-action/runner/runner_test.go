package runner

import (
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/logging"
	"testing"
)

func TestNewRunner(t *testing.T) {
	t.Run("check", func(t *testing.T) {
		cases := map[string]struct {
			Options  *Options
			Expected bool
		}{
		}

		for desc, c := range cases {
			t.Run(desc, func(t *testing.T) {
				spyLogger := logging.SpyLogger()

				run := NewRunner(checker.NewChecker(spyLogger))

				ok, err := run(c.Options)
				if err != nil {
					t.Errorf("want nil, got %#v", err)
					return
				}

				if ok != c.Expected {
					t.Errorf("want %t, got %t", c.Expected, ok)
				}
			})
		}
	})
	t.Run("check + junit", func(t *testing.T) {

	})
	t.Run("check + prcomment", func(t *testing.T) {

	})
	t.Run("check + autofix", func(t *testing.T) {

	})
	t.Run("check + junit + prcomment + autofix", func(t *testing.T) {

	})
}
