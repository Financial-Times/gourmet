package log

import "strings"

func LogLevelFromKeyword(level string) Level {
	logLevel := strings.ToLower(level)

	switch logLevel {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	}

	return InfoLevel
}
