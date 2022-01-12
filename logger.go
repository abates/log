package log

import (
	"io"
	"os"
	"sort"
	"sync"

	"github.com/abates/log/ansi"
)

type Logger interface {
	Printf(format string, v ...interface{}) TaskLogger
	Failf(format string, v ...interface{}) TaskLogger
	Successf(format string, v ...interface{}) TaskLogger
}

type printer interface {
	printf(block int, replace bool, mt MessageType, msg string)
}

type Cursor interface {
	Move(int) error
	Insert(int) error
	Hide() error
	Show() error
	Hidden() bool
	Clear() error
}

type insert struct {
	line  int
	lines int
}

type inserts []insert

func (ins inserts) Len() int           { return len(ins) }
func (ins inserts) Less(i, j int) bool { return ins[i].line < ins[j].line }
func (ins inserts) Swap(i, j int)      { ins[i], ins[j] = ins[j], ins[i] }

type logger struct {
	sync.Mutex

	inserts     inserts
	cursor      Cursor
	Format      Formatter
	out         io.Writer
	lines       int
	currentLine int
	end         int
}

func New(out io.Writer, formatters ...Formatter) Logger {
	l := &logger{
		out:    out,
		cursor: &ansi.Cursor{Writer: out},
	}

	if len(formatters) == 0 {
		l.Format = PlainFormatter()
	} else if len(formatters) == 1 {
		l.Format = formatters[0]
	} else {
		l.Format = FormatChain(formatters...)
	}

	return l
}

var std = New(os.Stderr)

func (l *logger) moveto(line int) int {
	moved := 0
	if l.currentLine < line {
		moved = line - l.currentLine
	} else if l.currentLine > line {
		moved = -(l.currentLine - line)
	}

	if moved != 0 {
		if !l.cursor.Hidden() {
			l.cursor.Hide()
		}
		l.cursor.Move(moved)
	}
	l.currentLine = line
	return moved
}

func (l *logger) printf(block int, replace bool, mt MessageType, msg string) {
	l.Lock()
	defer l.Unlock()
	if l.currentLine != l.end {
		l.moveto(l.end)
	}

	i, offset := l.getOffset(block)
	line := block + offset
	if replace {
		if i != nil && i.line == block {
			line -= i.lines
		}
		l.write(line, mt, msg)
		return
	}

	if line == l.end {
		if i != nil && i.line == block {
			line++
			i.lines++
		} else if i == nil {
			l.lines++
		}

		l.write(line, mt, msg)
		l.out.Write([]byte("\n"))
		l.end++
		l.currentLine++
		return
	}

	// insert
	if i == nil || i.line != block {
		l.inserts = append(l.inserts, insert{line: block, lines: 1})
		sort.Sort(l.inserts)
		line = block + offset + 1
	} else {
		line++
		i.lines++
	}

	if line == l.end {
		l.write(line, mt, msg)
		l.out.Write([]byte("\n"))
		l.end++
		l.currentLine++
	} else {
		l.out.Write([]byte("\n"))
		l.end++
		l.currentLine++

		l.moveto(line)
		l.cursor.Insert(1)
		l.currentLine = line
		l.write(line, mt, msg)
	}
}

func (l *logger) getOffset(line int) (i *insert, offset int) {
	offset = 0
	for index, insert := range l.inserts {
		if insert.line == line {
			return &l.inserts[index], offset + insert.lines
		} else if insert.line > line {
			if index > 0 {
				return &l.inserts[index-1], offset
			}
			return nil, offset
		}
		offset += insert.lines
	}
	return nil, offset
}

func (l *logger) write(line int, mt MessageType, msg string) Logger {
	if l.moveto(line) != 0 {
		l.cursor.Clear()
	}
	l.out.Write([]byte(msg))
	l.currentLine = line
	l.moveto(l.end)
	if l.cursor.Hidden() {
		l.cursor.Show()
	}
	return l
}

func (l *logger) nextLine() (line int) {
	l.Lock()
	line = l.lines
	l.Unlock()
	return line
}

func Printf(format string, v ...interface{}) TaskLogger {
	return std.Printf(format, v...)
}

func (l *logger) Printf(format string, v ...interface{}) TaskLogger {
	return l.taskLogger().Printf(format, v...)
}

func Failf(format string, v ...interface{}) TaskLogger {
	return std.Failf(format, v...)
}

func (l *logger) Failf(format string, v ...interface{}) TaskLogger {
	return l.taskLogger().Failf(format, v...)
}

func Successf(format string, v ...interface{}) TaskLogger {
	return std.Successf(format, v...)
}

func (l *logger) Successf(format string, v ...interface{}) TaskLogger {
	return l.taskLogger().Successf(format, v...)
}

func (l *logger) taskLogger() TaskLogger {
	return &taskLogger{printer: l, format: l.Format, line: l.nextLine()}
}
