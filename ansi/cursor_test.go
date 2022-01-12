package ansi

import (
	"bytes"
	"testing"
)

func TestCursor(t *testing.T) {
	tests := []struct {
		name string
		test func(*Cursor)
		want []byte
	}{
		{"move up", func(c *Cursor) { c.Move(-4) }, MoveUp(4)},
		{"move down", func(c *Cursor) { c.Move(3) }, MoveDown(3)},
		{"insert", func(c *Cursor) { c.Insert(-1) }, InsertLine(-1)},
		{"insert", func(c *Cursor) { c.Insert(2) }, InsertLine(2)},
		{"show", func(c *Cursor) { c.Show() }, ShowCursor.Write()},
		{"hide", func(c *Cursor) { c.Hide() }, HideCursor.Write()},
		{"clear", func(c *Cursor) { c.Clear() }, ClearLine()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			c := &Cursor{buf}
			test.test(c)
			got := buf.Bytes()

			if !bytes.Equal(test.want, got) {
				t.Errorf("Wanted %x got %x", test.want, got)
			}
		})
	}
}
