package options

import (
    "github.com/DeNA/unity-meta-check/util/typedpath"
    "testing"
)

func TestFakeRootDirValidator(t *testing.T) {
    cases := map[string]struct{
        Cwd typedpath.SlashPath
        RootDir typedpath.SlashPath
        Expected typedpath.SlashPath
    } {
        "already absolute": {
        	RootDir: "/already/absolute",
        	Expected: "/already/absolute",
        },
        "only .": {
            Cwd: "/path/to/cwd",
            RootDir: ".",
            Expected: "/path/to/cwd",
        },
        "relative start with ./": {
            Cwd: "/path/to/cwd",
            RootDir: "./rel",
            Expected: "/path/to/cwd/rel",
        },
        "relative not start with ./": {
            Cwd: "/path/to/cwd",
            RootDir: "rel",
            Expected: "/path/to/cwd/rel",
        },
    }

    for desc, c := range cases {
        t.Run(desc, func(t *testing.T) {
            validateRootDir := FakeRootDirValidator(c.Cwd.ToRaw())
        	actual, err := validateRootDir(c.RootDir.ToRaw())
        	if err != nil {
        	    t.Errorf("want nil, got %#v", err)
        	    return
            }

            if actual.ToSlash() != c.Expected {
                t.Errorf("want %q, got %q", c.Expected, actual.ToSlash())
            }
        })
    }
}
