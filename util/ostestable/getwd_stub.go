package ostestable

import "github.com/DeNA/unity-meta-check/util/typedpath"

//goland:noinspection GoUnusedExportedFunction
func StubGetwd(cwd typedpath.RawPath, err error) Getwd {
	return func() (typedpath.RawPath, error) {
		return cwd, err
	}
}
