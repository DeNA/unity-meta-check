package unity

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func ConstMetaNecessity(result bool) MetaNecessity {
	return func(typedpath.SlashPath) bool {
		return result
	}
}
