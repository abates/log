package ansi

import (
	"fmt"
)

type Sequence string

func (s Sequence) Write() []byte {
	return []byte(s)
}

func (s Sequence) Format(v interface{}) []byte {
	if i, ok := v.(int); ok && i < 0 {
		return []byte(s)
	}
	return []byte(fmt.Sprintf(string(s), v))
}

const (
	CSI Sequence = "\033["

	IL         = CSI + "%dL"  // Insert Line
	CPL        = CSI + "%dF"  // Cursor Previous Line
	CNL        = CSI + "%dE"  // Cusor Next Line
	EL         = CSI + "%dK"  // Erase Line
	HideCursor = CSI + "?25l" // Hide Cursor
	ShowCursor = CSI + "?25h" // Show Cursor
)

type Color func(string) string

var (
	Reset string = "\033[" + "0m"
	Red   Color  = genColor(31)
	Green Color  = genColor(32)
)

func genColor(c int) Color {
	color := fmt.Sprintf("%s%dm", CSI, c)
	return func(msg string) string {
		return fmt.Sprintf("%s%s%s", color, msg, Reset)
	}
}

func ClearLine() []byte {
	return EL.Format(2)
}

func InsertLine(num int) []byte {
	return IL.Format(num)
}

func MoveDown(num int) []byte {
	return CNL.Format(num)
}

func MoveUp(num int) []byte {
	return CPL.Format(num)
}
