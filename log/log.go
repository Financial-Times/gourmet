package log

type Level byte

const (
	TraceLevel Level = 1 << iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Event struct {
	fields map[string]interface{}
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
		fields: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (h *Event) Put(key string, val interface{}) {
	h.fields[key] = val
}

func (h *Event) Fields() map[string]interface{} {
	return h.fields
}

func (h *Event) clear() {
	h.fields = make(map[string]interface{})
}
