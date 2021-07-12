package log

type FluentLogger struct {
	*StructuredLogger
}

func NewFluentLogger(logger *StructuredLogger) *FluentLogger {
	return &FluentLogger{
		logger,
	}
}

func (f *FluentLogger) WithServiceName(val string) Logger  {
	return f.WithField("service_name", val)
}

func (f *FluentLogger) WithTransactionID(val string) Logger  {
	return f.WithField("transaction_id", val)
}

