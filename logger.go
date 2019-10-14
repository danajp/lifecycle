package lifecycle

// Logger defines behavior required by a logging package used by lifecycle phases
type Logger interface {
	Debug(msg string)
	Debugf(fmt string, v ...interface{})

	Info(msg string)
	Infof(fmt string, v ...interface{})

	Warn(msg string)
	Warnf(fmt string, v ...interface{})

	Error(msg string)
	Errorf(fmt string, v ...interface{})
}
