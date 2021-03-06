package interfaces

// Logger is the logger used in the app
type Logger interface {
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
}
