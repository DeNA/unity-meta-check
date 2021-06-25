package typedpath

import (
	"path"
	"path/filepath"
	"strings"
)

// SlashPath is slash separated path. Typically, Git uses the Slash Path.
type SlashPath string

func NewSlashPathUnsafe(path string) SlashPath {
	return SlashPath(path)
}

func SlashPathFromBaseName(baseName BaseName) SlashPath {
	return SlashPath(baseName)
}

func (s SlashPath) ToRaw() RawPath {
	return RawPath(filepath.FromSlash(string(s)))
}

func (s SlashPath) JoinSlashPath(other SlashPath) SlashPath {
	return SlashPath(path.Join(string(s), string(other)))
}

func (s SlashPath) JoinBaseName(other BaseName) SlashPath {
	return SlashPath(path.Join(string(s), string(other)))
}

func (s SlashPath) IsAbs() bool {
	return strings.HasPrefix(string(s), "/")
}

func (s SlashPath) Dir() SlashPath {
	return SlashPath(path.Dir(string(s)))
}

func (s SlashPath) Ext() string {
	return path.Ext(string(s))
}

func (s SlashPath) Split() (SlashPath, BaseName) {
	dirname, basename := path.Split(string(s))
	return SlashPath(dirname), BaseName(basename)
}
