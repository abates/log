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

func New(out io.Writer, formatters ...Formatter) *ProgressLogger {
	l := &ProgressLogger{
		out: out,
	}

	if len(formatters) == 0 {
		l.Format = PlainFormatter()
	} else if len(formatters) == 1 {
		l.Format = formatters[0]
	} else {
		l.Format = Formatters(formatters...)
	}

	return l
}

var std = New(os.Stderr)

func (l *ProgressLogger) moveto(line int) {
	l.out.Write(ansi.HideCursor.Write())
	if l.currentLine < line {
		l.out.Write(ansi.MoveDown(line - l.currentLine))
	} else if l.currentLine > line {
		l.out.Write(ansi.MoveUp(l.currentLine - line))
	}
	l.currentLine = line
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
	l.moveto(l.lines)
	l.out.Write(ansi.ShowCursor.Write())
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
	std.Printf(format, v...)
}

func (l *ProgressLogger) Printf(format string, v ...interface{}) {
	l.printf(l.nextLine(), InfoMessage, format, v...)
}

func Failf(format string, v ...interface{}) {
	std.Failf(format, v...)
}

func (l *ProgressLogger) Failf(format string, v ...interface{}) {
	l.printf(l.nextLine(), FailMessage, format, v...)
}

func Successf(format string, v ...interface{}) {
	std.Successf(format, v...)
}

func (l *ProgressLogger) Successf(format string, v ...interface{}) {
	l.printf(l.nextLine(), SuccessMessage, format, v...)
}

func LineLogger() Logger {
	return std.LineLogger()
}

func (l *ProgressLogger) LineLogger() Logger {
	return &lineLogger{Logger: l, line: l.nextLine()}
}
