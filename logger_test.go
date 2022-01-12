package log

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	logger := New(io.Discard)
	for n := 0; n < b.N; n++ {
		logger.Printf("This is a really really really long log message about absolutely nothing simply to benchmark the logger")
	}
}

type testCursor struct {
	commands []string
	hidden   bool
}

func (c *testCursor) Move(num int) error {
	c.commands = append(c.commands, fmt.Sprintf("move %d", num))
	return nil
}

func (c *testCursor) Insert(num int) error {
	c.commands = append(c.commands, fmt.Sprintf("insert %d", num))
	return nil
}

func (c *testCursor) Hidden() bool {
	return c.hidden
}

func (c *testCursor) Hide() error {
	c.hidden = true
	c.commands = append(c.commands, "hide")
	return nil
}

func (c *testCursor) Show() error {
	c.hidden = false
	c.commands = append(c.commands, "show")
	return nil
}

func (c *testCursor) Clear() error {
	c.commands = append(c.commands, "clear")
	return nil
}

func compare(t *testing.T, thing string, want, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %s %v got %v", thing, want, got)
	}
}

func TestLoggerMoveTo(t *testing.T) {
	tests := []struct {
		name        string
		currentLine int
		input       int
		want        []string
	}{
		{"move none", 0, 0, nil},
		{"move up", 25, 12, []string{"hide", "move -13"}},
		{"move down", 12, 25, []string{"hide", "move 13"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := New(io.Discard).(*logger)
			cursor := &testCursor{}
			l.cursor = cursor
			l.currentLine = test.currentLine
			l.moveto(test.input)
			compare(t, "commands", test.want, cursor.commands)
		})
	}

}

func checkState(t *testing.T, l *logger, lines int, currentLine int, inserts int) {
	t.Helper()
	if l.lines != lines {
		t.Errorf("Wanted %d lines got %d", lines, l.lines)
	}

	if l.currentLine != currentLine {
		t.Errorf("Wanted currentLine to be %d got %d", currentLine, l.currentLine)
	}

	if len(l.inserts) != inserts {
		t.Errorf("Wanted %d inserts got %d", inserts, len(l.inserts))
	}

}

func TestLoggerSimpleAppend(t *testing.T) {
	l := New(io.Discard).(*logger)
	cursor := &testCursor{}
	l.cursor = cursor

	l.Printf("One")
	checkState(t, l, 1, 1, 0)
	l.Printf("Two")
	checkState(t, l, 2, 2, 0)
	l.Printf("Three")
	checkState(t, l, 3, 3, 0)

	compare(t, "commands", []string(nil), cursor.commands)
}

func TestLoggerIntermediateAppend(t *testing.T) {
	l := New(io.Discard).(*logger)
	cursor := &testCursor{}
	l.cursor = cursor

	l.Printf("One")
	checkState(t, l, 1, 1, 0)
	l.Printf("Two")
	checkState(t, l, 2, 2, 0)
	ll := l.Printf("LL1")
	checkState(t, l, 3, 3, 0)
	ll.Printf("LL1: message 0")
	checkState(t, l, 3, 4, 1)

	if len(l.inserts) == 1 {
		if l.inserts[0].lines != 1 {
			t.Errorf("Expected 1 lines in the first insert got %d", l.inserts[0].lines)
		}

		if l.inserts[0].line != 2 {
			t.Errorf("Expected line to be 2 in first insert, got %d", l.inserts[0].line)
		}
	}

	compare(t, "commands", []string(nil), cursor.commands)
}

func TestLoggerAdvancedAppend(t *testing.T) {
	l := New(io.Discard).(*logger)
	cursor := &testCursor{}
	l.cursor = cursor

	ll1 := l.Printf("One")
	checkState(t, l, 1, 1, 0)
	l.Printf("Two")
	checkState(t, l, 2, 2, 0)
	ll2 := l.Printf("LL1")
	checkState(t, l, 3, 3, 0)
	ll2.Printf("LL1: message 0")
	checkState(t, l, 3, 4, 1)
	ll3 := l.Printf("Three")
	checkState(t, l, 4, 5, 1)
	ll2.Printf("LL1: message 1") // hide, move -2, insert 1, move 2, show
	checkState(t, l, 4, 6, 1)
	ll2.Printf("LL1: message 2") // hide, move -2, insert 1, move 2, show
	checkState(t, l, 4, 7, 1)
	ll4 := l.Printf("LL2")
	checkState(t, l, 5, 8, 1)
	ll4.Printf("LL2: message 0")
	checkState(t, l, 5, 9, 2)
	ll4.Printf("LL2: message 1")
	checkState(t, l, 5, 10, 2)
	ll3.Printf("Three: message 0") // hide, move -4, insert 1, move 4, show
	ll1.Printf("One: message 0")   // hide, move -11, insert 1, move 11, show

	if len(l.inserts) == 1 {
		if l.inserts[0].lines != 3 {
			t.Errorf("Expected 3 lines in the first insert got %d", l.inserts[0].lines)
		}

		if l.inserts[0].line != 2 {
			t.Errorf("Expected line to be 2 in first insert, got %d", l.inserts[0].line)
		}
	}
	want := []string{
		"hide", "move -2", "insert 1", "move 2", "show",
		"hide", "move -2", "insert 1", "move 2", "show",
		"hide", "move -4", "insert 1", "move 4", "show",
		"hide", "move -11", "insert 1", "move 11", "show",
	}

	compare(t, "commands", want, cursor.commands)
}
