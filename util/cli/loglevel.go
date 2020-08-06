package cli

import "github.com/DeNA/unity-meta-check/util/logging"

func GetLogLevel(debug, silent bool) logging.Severity {
	if debug {
		return logging.SeverityDebug
	}
	if silent {
		return logging.SeverityWarn
	}
	return logging.SeverityInfo
}
