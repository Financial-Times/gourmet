package log

import (
	"os"

	"github.com/rs/zerolog"
)

type StructuredLogger struct {
	logger zerolog.Logger
	logEvent *Event
	ref *Logger
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
		logEvent: NewLogEvent(),
	}
}

func (s *StructuredLogger) WithField(key string, val interface{}) *StructuredLogger {
	s.logEvent.Put(key, val)
	return s
}

func (s *StructuredLogger) Trace(message string, args ...interface{}) {
	s.log(s.logger.Trace(), message, args)
}

func (s *StructuredLogger) Debug(message string, args ...interface{}) {
	s.log(s.logger.Debug(), message, args)
}

func (s *StructuredLogger) Info(message string, args ...interface{}) {
	s.log(s.logger.Info(), message, args)
}

func (s *StructuredLogger) Warn(message string, args ...interface{}) {
	s.log(s.logger.Warn(), message, args)
}

func (s *StructuredLogger) Error(message string, args ...interface{}) {
	s.log(s.logger.Error(), message, args)
}

func (s *StructuredLogger) log(logEvent *zerolog.Event, message string, args []interface{}) {
	fields := s.logEvent.Fields()
	if len(fields) > 0 {
		logEvent = logEvent.Fields(fields)
	}
	logEvent.Fields(fields).Msgf(message, args)
	s.logEvent.Clear()
}
