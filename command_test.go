package command

import (
	"context"
	"testing"
)

func TestConcurrenceComE(t *testing.T) {
	err := ConcurrenceComNE(context.Background(),
		NewCmd("ls", 0,"-al"),
		NewCmd("lsss", 0,"-al"),
	)
	if err == nil {
		t.Error("err should not be empty")
	}
}

func TestConcurrenceComNE(t *testing.T) {
	err := ConcurrenceComNE(context.Background(),
		NewCmd("ls", 0,"-al"),
		NewCmd("lsss", 0,"-al"),
		)
	if err == nil {
		t.Error("err should not be empty")
	}
}