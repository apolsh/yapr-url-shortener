package logger

const ComponentKey = "component"

type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(err error)
	Fatal(err error)
}
