// Package logging defines the minimal interface that loggers must support to be used by pack.
package logging

import (
	"io"
)

const (
	InfoLevel  = "info"
	DebugLevel = "debug"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

// Logger defines behavior required by a logging package used by pack libraries
type Logger interface {
	Debug(msg string)
	Debugf(fmt string, v ...interface{})

	Info(msg string)
	Infof(fmt string, v ...interface{})

	Warn(msg string)
	Warnf(fmt string, v ...interface{})

	Error(msg string)
	Errorf(fmt string, v ...interface{})

	Writer() io.Writer

	WantLevel(level string)
}
