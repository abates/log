package log

import (
	"bytes"
	"fmt"
	golog "log"
	"strings"
	"sync"

	"github.com/abates/log/ansi"
)

type MessageType int

const (
	InfoMessage MessageType = iota
	SuccessMessage
	FailMessage
)

type Formatter func(mt MessageType, message string) string

func PlainFormatter() Formatter {
	return func(mt MessageType, message string) string { return message }
}

func ColorFormatter() Formatter {
	return func(mt MessageType, message string) string {
		switch mt {
		case SuccessMessage:
			message = ansi.Green(message)
		case FailMessage:
			message = ansi.Red(message)
		}
		return message
	}
}

func SuccessFormatter() Formatter {
	return func(mt MessageType, message string) string {
		switch mt {
		case SuccessMessage:
			message = fmt.Sprintf("%s %s", ansi.Green("✔"), message)
		case FailMessage:
			message = fmt.Sprintf("%s %s", ansi.Red("✕"), message)
		default:
			message = fmt.Sprintf("  %s", message)
		}
		return message
	}
}

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

func PrefixFormatter(prefix string, flag int) Formatter {
	buf := &bytes.Buffer{}
	mu := sync.Mutex{}
	logger := golog.New(buf, prefix, flag)
	return func(mt MessageType, message string) string {
		mu.Lock()
		logger.Print(message)
		message = strings.TrimSpace(string(buf.Bytes()))
		buf.Reset()
		mu.Unlock()
		return message
	}
}

func Formatters(formatters ...Formatter) Formatter {
	return func(mt MessageType, message string) string {
		for i := len(formatters) - 1; i >= 0; i-- {
			message = formatters[i](mt, message)
		}
		return message
	}
}
