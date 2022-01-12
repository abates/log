package log

import (
	"testing"
)

type testPrinter struct {
	gotStrings  []string
	gotTypes    []MessageType
	gotReplaces []bool
}

func (tl *testPrinter) printf(line int, replace bool, mt MessageType, msg string) {
	tl.gotReplaces = append(tl.gotReplaces, replace)
	tl.gotTypes = append(tl.gotTypes, mt)
	tl.gotStrings = append(tl.gotStrings, msg)
}

func TestTaskLogger(t *testing.T) {
	testPrinter := &testPrinter{}
	tl := &taskLogger{
		printer: testPrinter,
		format:  PlainFormatter(),
		line:    42,
	}
	tl.Successf("success")
	tl.Failf("fail")
	tl.Printf("print")
	tl.Setf(SuccessMessage, "set")

	wantStrings := []string{
		"success",
		"fail",
		"print",
		"set",
	}

	wantTypes := []MessageType{
		SuccessMessage,
		FailMessage,
		InfoMessage,
		SuccessMessage,
	}

	wantReplaces := []bool{
		false,
		false,
		false,
		true,
	}

	compare(t, "message types", wantTypes, testPrinter.gotTypes)
	compare(t, "strings", wantStrings, testPrinter.gotStrings)
	compare(t, "replaces", wantReplaces, testPrinter.gotReplaces)
}
