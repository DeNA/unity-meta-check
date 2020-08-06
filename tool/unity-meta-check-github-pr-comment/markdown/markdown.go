package markdown

import (
	"errors"
	"fmt"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"io"
	"strings"
)

func WriteMarkdown(result *checker.CheckResult, tmpl *l10n.Template, writer io.Writer) error {
	if result.Empty() {
		_, _ = fmt.Fprintln(writer, tmpl.SuccessMessage)
		return nil
	}

	_, _ = fmt.Fprintln(writer, tmpl.FailureMessage)
	_, _ = io.WriteString(writer, "\n")

	if err := WriteTableRow(writer, tmpl.StatusHeader, tmpl.FilePathHeader); err != nil {
		return err
	}

	if err := WriteTableSep(writer, 2); err != nil {
		return err
	}

	for _, missingMeta := range result.MissingMeta {
		if err := WriteTableRow(writer, tmpl.StatusMissing, FormatAsInlineCode(string(missingMeta))); err != nil {
			return err
		}
	}

	for _, danglingMeta := range result.DanglingMeta {
		if err := WriteTableRow(writer, tmpl.StatusDangling, FormatAsInlineCode(string(danglingMeta))); err != nil {
			return err
		}
	}

	return nil
}

func WriteTableSep(writer io.Writer, num int) error {
	if num == 0 {
		return errors.New("must include at least one column")
	}

	_, _ = io.WriteString(writer, "|")
	_, _ = io.WriteString(writer, strings.Repeat(":--|", num))
	_, _ = io.WriteString(writer, "\n")
	return nil
}

func WriteTableRow(writer io.Writer, cols ...string) error {
	if len(cols) == 0 {
		return errors.New("must include at least one column")
	}

	_, _ = io.WriteString(writer, "| ")
	_, _ = io.WriteString(writer, cols[0])
	for _, col := range cols[1:] {
		_, _ = io.WriteString(writer, " | ")
		_, _ = io.WriteString(writer, col)
	}
	_, _ = io.WriteString(writer, " |")
	_, _ = io.WriteString(writer, "\n")

	return nil
}

func FormatAsInlineCode(s string) string {
	return fmt.Sprintf("`%s`", s)
}
