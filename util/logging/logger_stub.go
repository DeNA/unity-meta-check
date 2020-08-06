package logging

import (
	"bytes"
	"fmt"
)

// WHY: It is an exposed API.
//goland:noinspection GoUnusedExportedFunction
func NullLogger() Logger {
	return nullLogger{}
}

type nullLogger struct{}

var _ Logger = nullLogger{}

func (nullLogger) Debug(string) {}

func (nullLogger) Info(string) {}

func (nullLogger) Warn(string) {}

func (nullLogger) Error(string) {}

func (nullLogger) Log(Severity, string) {}

func SpyLogger() *LoggerSpy {
	return &LoggerSpy{
		Logs: &bytes.Buffer{},
	}
}

type LoggerSpy struct {
	Logs *bytes.Buffer
}

var _ Logger = &LoggerSpy{}

func (s *LoggerSpy) Debug(message string) {
	s.Log(SeverityDebug, message)
}

func (s *LoggerSpy) Info(message string) {
	s.Log(SeverityInfo, message)
}

func (s *LoggerSpy) Warn(message string) {
	s.Log(SeverityWarn, message)
}

func (s *LoggerSpy) Error(message string) {
	s.Log(SeverityError, message)
}

func (s *LoggerSpy) Log(severity Severity, message string) {
	s.Logs.WriteString(fmt.Sprintf("%s: %s\n", severity.String(), message))
}
