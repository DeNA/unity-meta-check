package typedpath

import (
	"os"
	"path/filepath"
	"strings"
)

// NOTE: It is OS depended path. This style is for directly handling file systems.
type RawPath string

func NewRawPath(basenames ...BaseName) RawPath {
	elements := make([]string, len(basenames))
	for i, basename := range basenames {
		elements[i] = string(basename)
	}
	return RawPath(filepath.Join(elements...))
}

func NewRawPathUnsafe(path string) RawPath {
	return RawPath(path)
}

func Getwd() (RawPath, error) {
	result, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return RawPath(result), err
}

func (r RawPath) ToSlash() SlashPath {
	return SlashPath(filepath.ToSlash(string(r)))
}

func (r RawPath) JoinRawPath(other RawPath) RawPath {
	return RawPath(filepath.Join(string(r), string(other)))
}

func (r RawPath) JoinBaseName(other BaseName) RawPath {
	return RawPath(filepath.Join(string(r), string(other)))
}

func (r RawPath) IsAbs() bool {
	return filepath.IsAbs(string(r))
}

func (r RawPath) Rel(path RawPath) (RawPath, error) {
	result, err := filepath.Rel(string(r), string(path))
	if err != nil {
		return "", err
	}
	return RawPath(result), nil
}

func (r RawPath) Dir() RawPath {
	return RawPath(filepath.Dir(string(r)))
}

func (r RawPath) Ext() string {
	return filepath.Ext(string(r))
}

func (r RawPath) TrimLastSep() RawPath {
	return RawPath(strings.TrimRight(string(r), string(filepath.Separator)))
}
