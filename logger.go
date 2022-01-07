package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/abates/log/ansi"
)

type lineLogger struct {
	Logger
	line int
}

func (l *lineLogger) Printf(format string, v ...interface{}) {
	l.printf(l.line, InfoMessage, format, v...)
}

func (l *lineLogger) Failf(format string, v ...interface{}) {
	l.printf(l.line, FailMessage, format, v...)
}

func (l *lineLogger) Successf(format string, v ...interface{}) {
	l.printf(l.line, SuccessMessage, format, v...)
}

type Logger interface {
	Printf(format string, v ...interface{})
	Failf(format string, v ...interface{})
	Successf(format string, v ...interface{})

	printf(int, MessageType, string, ...interface{})
}

type ProgressLogger struct {
	sync.Mutex

	Format      Formatter
	out         io.Writer
	lines       int
	currentLine int
}

func New(out io.Writer) *ProgressLogger {
	l := &ProgressLogger{
		Format: PlainFormatter(),
		out:    out,
	}

	l.out.Write(ansi.HideCursor.Write())
	return l
}

var Std = New(os.Stderr)

func (l *ProgressLogger) moveto(line int) {
	if l.currentLine < line {
		l.out.Write(ansi.MoveDown(line - l.currentLine))
	} else if l.currentLine > line {
		l.out.Write(ansi.MoveUp(l.currentLine - line))
	}
}

func (l *ProgressLogger) printf(line int, mt MessageType, format string, v ...interface{}) {
	msg := l.Format(mt, fmt.Sprintf(format, v...))
	l.Lock()
	l.moveto(line)
	l.out.Write(ansi.ClearLine())
	l.out.Write([]byte(msg))
	l.out.Write([]byte("\n"))
	l.currentLine = line + 1
	l.Unlock()
}

func (l *ProgressLogger) nextLine() (line int) {
	l.Lock()
	if l.lines > 0 {
		l.moveto(l.lines - 1)
		l.out.Write([]byte("\n"))
	}
	line = l.lines
	l.lines++
	l.currentLine = line
	l.Unlock()
	return line
}

func Printf(format string, v ...interface{}) {
	Std.Printf(format, v...)
}

func (l *ProgressLogger) Printf(format string, v ...interface{}) {
	l.printf(l.nextLine(), InfoMessage, format, v...)
}

func Failf(format string, v ...interface{}) {
	Std.Failf(format, v...)
}

func (l *ProgressLogger) Failf(format string, v ...interface{}) {
	l.printf(l.nextLine(), FailMessage, format, v...)
}

func Successf(format string, v ...interface{}) {
	Std.Successf(format, v...)
}

func (l *ProgressLogger) Successf(format string, v ...interface{}) {
	l.printf(l.nextLine(), SuccessMessage, format, v...)
}

func LineLogger() Logger {
	return Std.LineLogger()
}

func Finish() {
	Std.Finish()
}

func (l *ProgressLogger) Finish() {
	if l.currentLine < l.lines {
		l.out.Write(ansi.MoveDown(l.lines - l.currentLine - 1))
	}
	l.out.Write([]byte("\n"))
	l.out.Write(ansi.ShowCursor.Write())
}

func (l *ProgressLogger) LineLogger() Logger {
	return &lineLogger{Logger: l, line: l.nextLine()}
}
