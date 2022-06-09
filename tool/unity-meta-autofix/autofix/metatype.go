package autofix

import (
	"fmt"
	"strings"

	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type MetaType string

const (
	MetaTypeDefaultImporterFolder MetaType = "MetaTypeDefaultImporterFolder"
	MetaTypeTextScriptImporter    MetaType = "MetaTypeTextScriptImporter"
	MetaTypeMonoImporter          MetaType = "MetaTypeMonoImporter"
)

type MetaTypeDetector func(missingMeta typedpath.RawPath) (MetaType, error)

func NewMetaTypeDetector(isDir ostestable.IsDir) MetaTypeDetector {
	return func(missingMeta typedpath.RawPath) (MetaType, error) {
		originalPath := unity.TrimMetaFromRaw(missingMeta)
		dir, err := isDir(originalPath)
		if err != nil {
			return "", err
		}

		ext := originalPath.Ext()
		if dir {
			if ext != "" {
				return "", fmt.Errorf("should not create .meta because the directory may require special .meta: %s", originalPath)
			}
			return MetaTypeDefaultImporterFolder, nil
		}

		switch strings.ToLower(ext) {
		case ".json", ".bytes", ".csv", ".pb", ".txt", ".xml", ".proto", ".md", ".asmdef":
			return MetaTypeTextScriptImporter, nil
		case ".cs":
			return MetaTypeMonoImporter, nil
		default:
			switch originalPath.Base() {
			case "LICENSE":
				return MetaTypeTextScriptImporter, nil
			default:
				return "", fmt.Errorf("should not create .meta because the extension is not supported now: %s", originalPath)
			}
		}
	}
}
