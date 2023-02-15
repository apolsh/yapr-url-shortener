package logger

// ComponentKey компонент приложения
const ComponentKey = "component"

// Interface интерфейс логера
type Interface interface {
	// Debug записать сообщение с уровнем debug
	Debug(message string, args ...interface{})
	// Info записать сообщение с уровнем info
	Info(message string, args ...interface{})
	// Warn записать сообщение с уровнем warn
	Warn(message string, args ...interface{})
	// Error записать сообщение с уровнем error
	Error(err error)
	// Fatal записать сообщение с уровнем fatal и завершить работу
	Fatal(err error)
}
