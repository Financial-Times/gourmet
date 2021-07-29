package gmlog

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

type Field func(*Event)

type Logger interface {
	Trace(message string, args ...Field)
	Debug(message string, args ...Field)
	Info(message string, args ...Field)
	Warn(message string, args ...Field)
	Error(message string, args ...Field)
}

func newLogEvent(opts ...Field) *Event {
	e := &Event{
		fields: cmap.New(),
	}
	// apply the list of options to Server
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (h *Event) put(key string, val interface{}) {
	h.fields.Set(key, val)
}

func (h *Event) Fields() map[string]interface{} {
	return h.fields.Items()
}

func (h *Event) clear() {
	h.fields = cmap.New()
}

