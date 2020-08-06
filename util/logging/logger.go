package logging

import (
	"fmt"
	"io"
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
		_, err := fmt.Fprintln(logger.writer, message)

		if err != nil {
			panic(err)
		}
	}
}
