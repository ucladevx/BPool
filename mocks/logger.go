package mocks

import (
	"github.com/ucladevx/BPool/utils/auth"
)

// Logger is a logger that does nothing
type Logger struct {
	auth.Logger
}

func (l Logger) Debug(msg string, args ...interface{}) {}
func (l Logger) Error(msg string, args ...interface{}) {}
func (l Logger) Info(msg string, args ...interface{})  {}
func (l Logger) Warn(msg string, args ...interface{})  {}
func (l Logger) Panic(msg string, args ...interface{}) {}
