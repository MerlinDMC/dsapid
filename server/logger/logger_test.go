package logger

import (
	"testing"
)

func TestLoggerDummy(t *testing.T) {
	SetName("testlogger")
	SetLevel(ERROR)

	Warn("this is a warning")
	Info("this is an info")
	Errorf("this is an %s", "Error")
}
