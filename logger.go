package bpool

import (
	"go.uber.org/zap"
)

// Logger is a BPool logger
type Logger struct {
	z *zap.SugaredLogger
}

// Info wraps the zap infow
func (l *Logger) Info(msg string, args ...interface{}) {
	l.z.Infow(msg, args...)
}

// Debug wraps the zap Debugw
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.z.Debugw(msg, args...)
}

// Warn wraps the zap Warnw
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.z.Warnw(msg, args...)
}

// Error wraps the zap Errorw
func (l *Logger) Error(msg string, args ...interface{}) {
	l.z.Errorw(msg, args...)
}

// Panic wraps the zap Panicw
func (l *Logger) Panic(msg string, args ...interface{}) {
	l.z.Panicw(msg, args...)
}

// NewBPoolLogger creates a new logger
func NewBPoolLogger(l *zap.SugaredLogger) *Logger {
	return &Logger{
		z: l,
	}
}
