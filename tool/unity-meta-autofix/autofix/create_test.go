package autofix

import (
	"github.com/DeNA/unity-meta-check/unity/meta"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewMetaCreator(t *testing.T) {
	cases := []MetaType{
		MetaTypeDefaultImporterFolder,
		MetaTypeTextScriptImporter,
	}

	for _, metaType := range cases {
		t.Run(string(metaType), func(t *testing.T) {
			spyLogger := logging.SpyLogger()

			workDir, err := ioutil.TempDir(os.TempDir(), "")
			if err != nil {
				panic(err.Error())
			}
			missingMeta := typedpath.RawPath(filepath.Join(workDir, "Missing.meta"))

			createMeta := NewMetaCreator(false, meta.StubGUIDGen(meta.AnyGUID(), nil), spyLogger)
			if err := createMeta(MetaTypeDefaultImporterFolder, missingMeta); err != nil {
				t.Log(spyLogger.Logs.String())
				t.Errorf("want nil, got %#v", err)
				return
			}

			if _, err := os.Stat(string(missingMeta)); err != nil {
				t.Log(spyLogger.Logs.String())
				t.Errorf("want nil, got %#v", err)
				return
			}

			_ = os.RemoveAll(workDir)
		})
	}
}
