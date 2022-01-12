package ansi

import "testing"

func TestColor(t *testing.T) {
	c := genColor(31)
	want := "\033[31mHello World\033[0m"
	got := c("Hello World")

	if want != got {
		t.Errorf("Wanted %q got %q", want, got)
	}
}
