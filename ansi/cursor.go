package ansi

import "io"

type Cursor struct {
	io.Writer
	hidden bool
}

func (c *Cursor) write(data []byte) error {
	_, err := c.Write(data)
	return err
}

func (c *Cursor) Move(num int) error {
	if num < 0 {
		return c.write(MoveUp(-1 * num))
	}
	return c.write(MoveDown(num))
}

func (c *Cursor) Insert(num int) error {
	return c.write(InsertLine(num))
}

func (c *Cursor) Hide() error {
	c.hidden = true
	return c.write(HideCursor.Write())
}

func (c *Cursor) Hidden() bool {
	return c.hidden
}

func (c *Cursor) Show() error {
	c.hidden = false
	return c.write(ShowCursor.Write())
}

func (c *Cursor) Clear() error {
	return c.write(ClearLine())
}
