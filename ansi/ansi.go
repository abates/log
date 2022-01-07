package ansi

import (
	"fmt"
	"strconv"
)

type Sequence string

func (s Sequence) Write() []byte {
	return []byte(s)
}

func (s Sequence) Format(v int) []byte {
	if v < 0 {
		return []byte(fmt.Sprintf(string(s), ""))
	}
	return []byte(fmt.Sprintf(string(s), strconv.Itoa(v)))
}

const (
	CSI Sequence = "\033["

	CPL        = CSI + "%sF"
	CNL        = CSI + "%sE"
	EL         = CSI + "%sK"
	HideCursor = CSI + "?25l"
	ShowCursor = CSI + "?25h"
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

func MoveDown(num int) []byte {
	return CNL.Format(num)
}

func MoveUp(num int) []byte {
	return CPL.Format(num)
}
