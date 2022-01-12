package log

import (
	"fmt"
	"strings"

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
	return func(mt MessageType, message string) string {
		return message
	}
}

func Indent(i int) Formatter {
	indent := strings.Repeat(" ", i)
	return func(mt MessageType, message string) string {
		return fmt.Sprintf("%s%s", indent, message)
	}
}

func Colorize() Formatter {
	splitIndex := func(msg string) (string, string) {
		if index := strings.Index(msg, ":"); index > -1 && index < len(msg)-1 {
			return msg[0 : index+1], msg[index+1:]
		}
		return "", msg
	}

	return func(mt MessageType, message string) string {
		switch mt {
		case SuccessMessage:
			bef, aft := splitIndex(message)
			message = bef + ansi.Green(aft)
		case FailMessage:
			bef, aft := splitIndex(message)
			message = bef + ansi.Red(aft)
		}
		return message
	}
}

func Annotate() Formatter {
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

func FormatChain(formatters ...Formatter) Formatter {
	return func(mt MessageType, message string) string {
		for i := len(formatters) - 1; i >= 0; i-- {
			message = formatters[i](mt, message)
		}
		return message
	}
}
