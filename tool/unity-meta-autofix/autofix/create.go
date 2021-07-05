package autofix

import (
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/meta"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
)

type MetaCreator func(metaType MetaType, originalPath typedpath.RawPath) error

func NewMetaCreator(dryRun bool, guidGen meta.GUIDGen, logger logging.Logger) MetaCreator {
	return func(metaType MetaType, missingMeta typedpath.RawPath) error {
		guid, err := guidGen()
		if err != nil {
			return err
		}

		var metaGen meta.Gen
		switch metaType {
		case MetaTypeDefaultImporterFolder:
			metaGen = meta.DefaultImporterFolderGen{GUID: guid}
		case MetaTypeTextScriptImporter:
			metaGen = meta.TextScriptImporterGen{GUID: guid}
		case MetaTypeMonoImporter:
			metaGen = meta.MonoImporterGen{GUID: guid}
		default:
			return fmt.Errorf("unsupported meta type: %q", metaType)
		}

		_, err = os.Stat(string(missingMeta))
		if err == nil {
			return fmt.Errorf("file exists: %s", missingMeta)
		}
		if !os.IsNotExist(err) {
			return err
		}

		if dryRun {
			return createMetaDryRun(missingMeta, metaGen, logger)
		}
		return createMeta(missingMeta, metaGen)
	}
}

func createMetaDryRun(missingMeta typedpath.RawPath, metaGen meta.Gen, logger logging.Logger) error {
	buf := &bytes.Buffer{}
	if _, err := metaGen.WriteTo(buf); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("would write to %q:\n%s", missingMeta, buf.String()))
	return nil
}

func createMeta(missingMeta typedpath.RawPath, metaGen meta.Gen) error {
	file, err := os.OpenFile(string(missingMeta), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer func(){ _ = file.Close() }()

	_, err = metaGen.WriteTo(file)
	if err != nil {
		return err
	}
	return nil
}
