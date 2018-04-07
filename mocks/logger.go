package mocks

// Logger is a logger that does nothing
type Logger struct{}

// Debug is a mock for debug
func (l Logger) Debug(msg string, args ...interface{}) {}

// Error is a mock for Error
func (l Logger) Error(msg string, args ...interface{}) {}

// Info is a mock for info
func (l Logger) Info(msg string, args ...interface{}) {}

// Warn is a mock for Warn
func (l Logger) Warn(msg string, args ...interface{}) {}

// Panic is a mock for Panic
func (l Logger) Panic(msg string, args ...interface{}) {}
