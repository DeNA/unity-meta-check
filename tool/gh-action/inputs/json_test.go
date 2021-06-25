package inputs

import (
	"encoding/json"
	"github.com/DeNA/unity-meta-check/tool/gh-action/action-yaml-gen/yaml"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
)

func TestExample(t *testing.T) {
	exampleBytes, err := os.ReadFile("testdata/inputs-example.json")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	var in Inputs
	if err := json.Unmarshal(exampleBytes, &in); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := Inputs{
		LogLevel:                   "DEBUG",
		TargetPath:                 "/home/runner/work/unity-meta-check-playground/unity-meta-check-playground",
		TargetType:                 "auto-detect",
		IgnoreDangling:             false,
		IgnoreCase:                 false,
		IgnoreSubmodulesAndNested:  false,
		IgnoredFilePath:            ".meta-check-ignore",
		EnableAutofix:              true,
		CommaSeparatedAutofixGlobs: ".",
		EnableJUnit:                true,
		JUnitXMLPath:               "junit.xml",
		EnablePRComment:            true,
		PRCommentTmplFilePath:      "",
		PRCommentLang:              "ja",
		PRCommentSendSuccess:       true,
	}

	if !reflect.DeepEqual(expected, in) {
		t.Error(cmp.Diff(expected, in))
	}
}

func TestInputDefsCoverInput(t *testing.T) {
	inputJSONFieldNameMap, missingJSONTags := buildInputJSONFieldNameMap()
	for _, missingJSONTag := range missingJSONTags {
		t.Errorf("missing json tag: %q", missingJSONTag)
	}

	inputDefNameMap := buildInputDefNameMap()

	for inputJSONField := range inputJSONFieldNameMap {
		if _, ok := inputDefNameMap[yaml.Name(inputJSONField)]; !ok {
			t.Errorf("missing input field: %q", inputJSONField)
		}
	}

	for inputDefName := range inputDefNameMap {
		if _, ok := inputJSONFieldNameMap[jsonFieldName(inputDefName)]; !ok {
			t.Errorf("extra input definition: %q", inputDefName)
		}
	}
}

type jsonFieldName string
type missingFieldName string

func buildInputJSONFieldNameMap() (map[jsonFieldName]struct{}, []missingFieldName) {
	var i Inputs
	it := reflect.TypeOf(i)
	nf := it.NumField()

	res := make(map[jsonFieldName]struct{}, nf)
	missing := make([]missingFieldName, 0, nf)

	for i := 0; i < nf; i++ {
		field := it.Field(i)

		jsonTag, ok := field.Tag.Lookup("json")
		if !ok {
			missing = append(missing, missingFieldName(field.Name))
		}

		tagTokens := strings.SplitN(jsonTag, ",", 2)
		jsonFieldName := jsonFieldName(tagTokens[0])
		res[jsonFieldName] = struct{}{}
	}

	return res, missing
}

func buildInputDefNameMap() map[yaml.Name]struct{} {
	res := make(map[yaml.Name]struct{}, len(yaml.BuildMetadata().Inputs))

	for _, inputDef := range yaml.BuildMetadata().Inputs {
		res[inputDef.Name()] = struct{}{}
	}

	return res
}

func TestStringifyBool_UnmarshalJSON(t *testing.T) {
	cases := map[string]struct {
		JSON     string
		Expected StringifyBool
	}{
		"true": {
			JSON:     `"true"`,
			Expected: true,
		},
		"false": {
			JSON:     `"false"`,
			Expected: false,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			var actual StringifyBool

			if err := json.Unmarshal([]byte(c.JSON), &actual); err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if actual != c.Expected {
				t.Errorf("want %t, got %t", c.Expected, actual)
			}
		})
	}
}

func TestStringifyInt_UnmarshalJSON(t *testing.T) {
	cases := map[string]struct {
		JSON     string
		Expected StringifyInt
	}{
		"0 (boundary)": {
			JSON:     `"0"`,
			Expected: 0,
		},
		"1 (positive representative)": {
			JSON:     `"1"`,
			Expected: 1,
		},
		"-1 (negative representative)": {
			JSON:     `"-1"`,
			Expected: -1,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			var actual StringifyInt

			if err := json.Unmarshal([]byte(c.JSON), &actual); err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if actual != c.Expected {
				t.Errorf("want %d, got %d", c.Expected, actual)
			}
		})
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	if err := quick.Check(func(i1 Inputs) bool {
		j, err := json.Marshal(i1)
		if err != nil {
			t.Log(err)
			return false
		}

		var i2 Inputs
		if err := json.Unmarshal(j, &i2); err != nil {
			t.Log(err)
			return false
		}

		return reflect.DeepEqual(i1, i2)
	}, nil); err != nil {
		t.Error(err)
	}
}

func TestValidateTargetType(t *testing.T) {
	cases := []struct {
		TargetType string
		Expected   TargetType
	}{
		{
			TargetType: "auto-detect",
			Expected:   TargetTypeAuto,
		},
		{
			TargetType: "upm-package",
			Expected:   TargetTypeUpmPackage,
		},
		{
			TargetType: "unity-project",
			Expected:   TargetTypeUnityProj,
		},
		{
			TargetType: "unity-project-sub-dir",
			Expected:   TargetTypeUnityProjSubDir,
		},
	}

	for _, c := range cases {
		t.Run(c.TargetType, func(t *testing.T) {
			actual, err := ValidateTargetType(c.TargetType)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if actual != c.Expected {
				t.Errorf("want %q, got %q", c.Expected, actual)
			}
		})
	}
}
