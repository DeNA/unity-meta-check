package logging

import (
	"fmt"
	"io"
	"strings"
)

type Severity int

const (
	SeverityDebug Severity = iota
	SeverityInfo  Severity = iota
	SeverityWarn  Severity = iota
	SeverityError Severity = iota
)

func (s Severity) String() string {
	switch s {
	case SeverityDebug:
		return "DEBUG"
	case SeverityInfo:
		return "INFO"
	case SeverityWarn:
		return "WARN"
	case SeverityError:
		return "ERROR"
	default:
		panic("unreachable")
	}
}

func ParseSeverity(s string) (Severity, error) {
	switch s {
	case "DEBUG":
		return SeverityDebug, nil
	case "INFO":
		return SeverityInfo, nil
	case "WARN":
		return SeverityWarn, nil
	case "ERROR":
		return SeverityError, nil
	default:
		return 0, fmt.Errorf("unknown severity: %q", s)
	}
}

// MustParseSeverity return the severity if given "DEBUG"/"INFO"/"WARN"/"ERROR". otherwise return "DEBUG" to fallback.
func MustParseSeverity(s string) Severity {
	v, err := ParseSeverity(s)
	if err != nil {
		// NOTE: Fallback
		return SeverityDebug
	}
	return v
}

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Log(Severity, string)
}

func NewLogger(severity Severity, writer io.Writer) Logger {
	return &severityLogger{
		severity: severity,
		writer:   writer,
	}
}

type severityLogger struct {
	severity Severity
	writer   io.Writer
}

func (logger *severityLogger) Debug(message string) {
	logger.Log(SeverityDebug, message)
}

func (logger *severityLogger) Info(message string) {
	logger.Log(SeverityInfo, message)
}

func (logger *severityLogger) Warn(message string) {
	logger.Log(SeverityWarn, message)
}

func (logger *severityLogger) Error(message string) {
	logger.Log(SeverityError, message)
}

func (logger *severityLogger) Log(severity Severity, message string) {
	if logger.severity <= severity {
		_, _ = fmt.Fprintf(logger.writer, "%s: %s\n", strings.ToLower(severity.String()), message)
	}
}
