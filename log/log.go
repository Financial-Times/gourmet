package log


type Level byte

const (
	Debug Level = 1 << iota
	Info
	Warn
	Error
)

type Logger interface {
	Debug(keyVals ...interface{})
	Info(keyVals ...interface{})
	Warn(keyVals ...interface{})
	Error(keyVals ...interface{})
}
