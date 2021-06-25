package log

import (
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type StructuredLogger struct {
	kitlog.Logger
	keyConf *KeyNamesConfig
}

func NewStructuredLogger(logLevel Level, customLogKeyConf ...*KeyNamesConfig) StructuredLogger {
	var rawLogLevel level.Option
	switch logLevel {
	case Debug:
		rawLogLevel = level.AllowDebug()
	case Warn:
		rawLogLevel = level.AllowWarn()
	case Error:
		rawLogLevel = level.AllowError()
	default:
		rawLogLevel = level.AllowInfo()
	}

	logger := kitlog.NewJSONLogger(os.Stdout)
	logger = level.NewFilter(logger, rawLogLevel)

	keyConf := NewDefaultKeyNamesConfig()
	if len(customLogKeyConf) > 0 {
		keyConf = NewKeyNamesConfig(customLogKeyConf[0])
	}
	logger = kitlog.With(logger, keyConf.KeyTime, kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, keyConf.KeyCaller, kitlog.DefaultCaller)
	return StructuredLogger{
		keyConf: keyConf,
	}
}

func (s *StructuredLogger) Debug(keyVals ...interface{}) {
	keyVals = s.appendMessageKeyIfNeeded(keyVals)
	_ = level.Debug(s).Log(keyVals)
}

func (s *StructuredLogger) Info(keyVals ...interface{}) {
	keyVals = s.appendMessageKeyIfNeeded(keyVals)
	_ = level.Info(s).Log(keyVals)
}

func (s *StructuredLogger) Warn(keyVals ...interface{}) {
	keyVals = s.appendMessageKeyIfNeeded(keyVals)
	_ = level.Warn(s).Log(keyVals)
}

func (s *StructuredLogger) Error(keyVals ...interface{}) {
	keyVals = s.appendMessageKeyIfNeeded(keyVals)
	_ = level.Error(s).Log(keyVals)
}

func (s *StructuredLogger) appendMessageKeyIfNeeded(keyVals ...interface{}) []interface{} {
	if len(keyVals) == 1 {
		keyVals = append(keyVals, 0)
		copy(keyVals[1:], keyVals)
		keyVals[0] = s.keyConf.KeyMessage
	}
	return keyVals
}
