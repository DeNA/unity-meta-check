package inputs

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Inputs struct {
	LogLevel string `json:"log_level"`

	TargetType                string        `json:"target_type"`
	IgnoreDangling            StringifyBool `json:"ignore_dangling"`
	IgnoreCase                StringifyBool `json:"ignore_case"`
	IgnoreSubmodulesAndNested StringifyBool `json:"ignore_submodules_and_nested"`
	IgnoredFilePath           string        `json:"ignored_file_path"`

	EnableAutofix   StringifyBool `json:"enable_autofix"`
	AutofixMissing  StringifyBool `json:"autofix_missing"`
	AutofixDangling StringifyBool `json:"autofix_dangling"`
	AutofixGlobs    StringifyBool `json:"autofix_globs"`

	EnableJUnit  StringifyBool `json:"enable_junit"`
	JUnitXMLPath string        `json:"junit_xml_path"`

	EnablePRComment        StringifyBool `json:"enable_pr_comment"`
	PRCommentPRNumber      StringifyInt  `json:"pr_comment_pull"`
	PRCommentTmplFilePath  string        `json:"pr_comment_tmpl_file"`
	PRCommentLang          string        `json:"pr_comment_lang"`
	PRCommentOwner         string        `json:"pr_comment_owner"`
	PRCommentRepo          string        `json:"pr_comment_repo"`
	PRCommentAPIEndpoint   string        `json:"pr_comment_api_endpoint"`
	PRCommentNoSendSuccess StringifyBool `json:"pr_comment_no_send_success"`
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
		return []byte("true"), nil
	}
	return []byte("false"), nil
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
	return []byte(strconv.Itoa(int(s))), nil
}
