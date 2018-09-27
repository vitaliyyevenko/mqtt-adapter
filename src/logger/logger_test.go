package logger

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	if Log == nil{
		t.Error("logger shouldn't be <nil>")
	}
}
