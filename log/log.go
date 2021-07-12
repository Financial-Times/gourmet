package log

import cmap "github.com/orcaman/concurrent-map"

type Level byte

const (
	TraceLevel Level = 1 << iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Event struct {
	fields cmap.ConcurrentMap
}

func NewLogEvent() *Event {
	return &Event{
		fields: cmap.New(),
	}
}

func (h *Event) Put(key string, val interface{}) {
	h.fields.Set(key, val)
	return h
}

func (h *Event) Fields() map[string]interface{} {
	return h.fields.Items()
}

func (h *Event) Clear() {
	h.fields = cmap.New()
}

type Logger interface {
	WithField(key string, val interface{}) *Logger
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

