package log

func WithField(key string, val interface{}) Field {
	return func(e *Event) {
		e.Put(key, val)
	}
}

func WithTransactionID(tid string) Field {
	return WithField("transaction_id", tid)
}

func WithServiceName(name string) Field {
	return WithField("service_name", name)
}

func WithError(err error) Field {
	return WithField("error", err.Error())
}
