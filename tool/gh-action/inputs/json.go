package inputs

import (
	"encoding/json"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"strconv"
)

type TargetType string

const (
	TargetTypeAuto            TargetType = "auto-detect"
	TargetTypeUnityProj       TargetType = "unity-project"
	TargetTypeUnityProjSubDir TargetType = "unity-project-sub-dir"
	TargetTypeUpmPackage      TargetType = "upm-package"
)

func ValidateTargetType(unsafeTargetType string) (TargetType, error) {
	switch TargetType(unsafeTargetType) {
	case TargetTypeAuto:
		return TargetTypeAuto, nil
	case TargetTypeUnityProj:
		return TargetTypeUnityProj, nil
	case TargetTypeUnityProjSubDir:
		return TargetTypeUnityProjSubDir, nil
	case TargetTypeUpmPackage:
		return TargetTypeUpmPackage, nil
	default:
		return "", fmt.Errorf("unknown target type: %q", unsafeTargetType)
	}
}

type Inputs struct {
	LogLevel string `json:"log_level"`

	TargetType                string              `json:"target_type"`
	IgnoreDangling            StringifyBool       `json:"ignore_dangling"`
	IgnoreCase                StringifyBool       `json:"ignore_case"`
	IgnoreSubmodulesAndNested StringifyBool       `json:"ignore_submodules_and_nested"`
	IgnoredFilePath           typedpath.SlashPath `json:"ignored_file_path"`

	EnableAutofix StringifyBool `json:"enable_autofix"`
	AutofixGlobs  []string      `json:"autofix_globs"`

	EnableJUnit  StringifyBool       `json:"enable_junit"`
	JUnitXMLPath typedpath.SlashPath `json:"junit_xml_path"`

	EnablePRComment       StringifyBool       `json:"enable_pr_comment"`
	PRCommentPRNumber     StringifyInt        `json:"pr_comment_pull"`
	PRCommentTmplFilePath typedpath.SlashPath `json:"pr_comment_tmpl_file"`
	PRCommentLang         string              `json:"pr_comment_lang"`
	PRCommentOwner        string              `json:"pr_comment_owner"`
	PRCommentRepo         string              `json:"pr_comment_repo"`
	PRCommentAPIEndpoint  string              `json:"pr_comment_api_endpoint"`
	PRCommentSendSuccess  StringifyBool       `json:"pr_comment_send_success"`
}

type StringifyBool bool

func (s *StringifyBool) UnmarshalJSON(bytes []byte) error {
	switch string(bytes) {
	case `"true"`:
		*s = true
		return nil
	case `"false"`:
		*s = false
		return nil
	default:
		return fmt.Errorf("illegal StringifyBool: %q", bytes)
	}
}

func (s StringifyBool) MarshalJSON() ([]byte, error) {
	if s {
		return []byte(`"true"`), nil
	}
	return []byte(`"false"`), nil
}

type StringifyInt int

func (s *StringifyInt) UnmarshalJSON(bytes []byte) error {
	var s2 string
	if err := json.Unmarshal(bytes, &s2); err != nil {
		return fmt.Errorf("illegal StringifyInt: %q", bytes)
	}

	i, err := strconv.Atoi(s2)
	if err != nil {
		return fmt.Errorf("illegal StringifyInt: %q", bytes)
	}
	*s = StringifyInt(i)
	return nil
}

func (s StringifyInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%q`, strconv.Itoa(int(s)))), nil
}
