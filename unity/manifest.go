package unity

import (
	"encoding/json"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io/ioutil"
)

type ManifestJson struct {
	Dependencies map[string]string `json:"dependencies"`
}

var ManifestBasename typedpath.BaseName = "manifest.json"

func ReadManifest(path typedpath.RawPath) (*ManifestJson, error) {
	bytes, err := ioutil.ReadFile(string(path))
	if err != nil {
		return nil, err
	}
	return parseManifestJson(bytes)
}

func parseManifestJson(bytes []byte) (*ManifestJson, error) {
	var manifestJson ManifestJson
	if err := json.Unmarshal(bytes, &manifestJson); err != nil {
		return nil, err
	}
	return &manifestJson, nil
}
