package l10n

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io"
	"io/ioutil"
	"os"
)

type Template struct {
	SuccessMessage string `json:"success"`
	FailureMessage string `json:"failure"`
	StatusHeader   string `json:"header_status"`
	FilePathHeader string `json:"header_file_path"`
	StatusMissing  string `json:"status_missing"`
	StatusDangling string `json:"status_dangling"`
}

var En = Template{
	SuccessMessage: "No missing/dangling .meta found. Perfect!",
	FailureMessage: `Some missing or dangling .meta found. Fix commits are needed.`,
	StatusHeader:   "Status",
	FilePathHeader: "File",
	StatusMissing:  "Not committed",
	StatusDangling: "Not removed",
}

var Ja = Template{
	SuccessMessage: "commit忘れ・消し忘れの .meta はありませんでした。素晴らしい！",
	FailureMessage: "commit忘れ・消し忘れの .meta が見つかりました。修正コミットが必要です。",
	StatusHeader:   "状態",
	FilePathHeader: "ファイル",
	StatusMissing:  "commit されていない",
	StatusDangling: "消されていない",
}

func GetTemplate(lang Lang) (*Template, error) {
	switch lang {
	case LangEn:
		return &En, nil
	case LangJa:
		return &Ja, nil
	default:
		return nil, fmt.Errorf("unsupported lang: %s", lang)
	}
}

func ReadTemplateFile(path typedpath.RawPath) (*Template, error) {
	file, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}
	defer func(){ _ = file.Close() }()

	return ReadTemplate(file)
}

func ReadTemplate(reader io.Reader) (*Template, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var tmpl Template
	if err := json.Unmarshal(bytes, &tmpl); err != nil {
		return nil, err
	}

	return &tmpl, nil
}

func WriteTemplateExample(writer io.Writer) {
	bytes, err := json.MarshalIndent(En, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	_, _ = writer.Write(bytes)
}

func ValidateTemplate(tmpl *Template) error {
	if tmpl.StatusHeader == "" {
		return errors.New(`empty "header_status"`)
	}
	if tmpl.FilePathHeader == "" {
		return errors.New(`empty "header_file_path"`)
	}
	if tmpl.StatusMissing == "" {
		return errors.New(`empty "status_missing"`)
	}
	if tmpl.StatusDangling == "" {
		return errors.New(`empty "status_dangling"`)
	}
	return nil
}
