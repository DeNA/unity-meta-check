package ignore

import (
	"bufio"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io"
	"os"
	"strings"
)

func ReadFile(path typedpath.RawPath) ([]globs.Glob, error) {
	file, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}
	defer func(){ _ = file.Close() }()
	return Read(file)
}

func Read(reader io.Reader) ([]globs.Glob, error) {
	result := make([]globs.Glob, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		elementsRaw := strings.SplitN(line, "#", 2)
		globPattern := strings.Trim(strings.TrimSpace(elementsRaw[0]), "/")
		if globPattern != "" {
			result = append(result, globs.Glob(globPattern))
		}
	}
	return result, nil
}
