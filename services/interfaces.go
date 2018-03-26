package services

type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Panic(args ...interface{})
}
