package gmlog

import (
	"os"

	"github.com/rs/zerolog"
)

type StructuredLogger struct {
	logger zerolog.Logger
}

func NewStructuredLogger(logLevel Level) *StructuredLogger {
	var rawLogLevel zerolog.Level
	switch logLevel {
	case TraceLevel:
		rawLogLevel = zerolog.TraceLevel
	case DebugLevel:
		rawLogLevel = zerolog.DebugLevel
	case WarnLevel:
		rawLogLevel = zerolog.WarnLevel
	case ErrorLevel:
		rawLogLevel = zerolog.ErrorLevel
	default:
		rawLogLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).Level(rawLogLevel).With().
		Timestamp().Caller().
		Logger()

	return &StructuredLogger{
		logger: logger,
	}
}

func (s *StructuredLogger) Trace(message string, args ...Field) {
	logInternal(s.logger.Trace(), message, args...)
}

func (s *StructuredLogger) Debug(message string, args ...Field) {
	logInternal(s.logger.Debug(), message, args...)
}

func (s *StructuredLogger) Info(message string, args ...Field) {
	logInternal(s.logger.Info(), message, args...)
}

func (s *StructuredLogger) Warn(message string, args ...Field) {
	logInternal(s.logger.Info(), message, args...)
}

func (s *StructuredLogger) Error(message string, args ...Field) {
	logInternal(s.logger.Error(), message, args...)
}

func logInternal(logEvent *zerolog.Event, message string, args ...Field) {
	e := newLogEvent(args...)
	fields := e.Fields()
	if len(fields) > 0 {
		logEvent = logEvent.Fields(fields)
	}
	logEvent.Fields(fields).Msgf(message)
	e.clear()
}
