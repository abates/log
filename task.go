package log

import "fmt"

type TaskLogger interface {
	Logger
	Setf(mt MessageType, format string, v ...interface{}) TaskLogger
	SetFormatter(Formatter)
	Formatter() Formatter
}

type taskLogger struct {
	printer
	format Formatter
	line   int
}

func (tl *taskLogger) SetFormatter(format Formatter) { tl.format = format }
func (tl *taskLogger) Formatter() Formatter          { return tl.format }

func (tl *taskLogger) printf(replace bool, mt MessageType, format string, v ...interface{}) TaskLogger {
	msg := tl.format(mt, fmt.Sprintf(format, v...))
	tl.printer.printf(tl.line, replace, mt, msg)
	return tl
}

func (tl *taskLogger) Setf(mt MessageType, format string, v ...interface{}) TaskLogger {
	return tl.printf(true, mt, format, v...)
}

func (tl *taskLogger) Printf(format string, v ...interface{}) TaskLogger {
	return tl.printf(false, InfoMessage, format, v...)
}

func (tl *taskLogger) Failf(format string, v ...interface{}) TaskLogger {
	return tl.printf(false, FailMessage, format, v...)
}

func (tl *taskLogger) Successf(format string, v ...interface{}) TaskLogger {
	return tl.printf(false, SuccessMessage, format, v...)
}
