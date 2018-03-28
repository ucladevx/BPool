package mocks

import (
	"github.com/ucladevx/BPool/utils/auth"
)

// Logger is a logger that does nothing
type Logger struct {
	auth.Logger
}

func (l Logger) Debug(args ...interface{}) {}
func (l Logger) Error(args ...interface{}) {}
func (l Logger) Info(args ...interface{})  {}
func (l Logger) Warn(args ...interface{})  {}
func (l Logger) Panic(args ...interface{}) {}
