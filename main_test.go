package main

import (
	"bytes"
	"testing"
)

func TestRun_MissingFile(t *testing.T) {
	err := run([]string{"check", "does-not-exist.html"}, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_NoCommand(t *testing.T) {
	err := run([]string{}, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error when no command given")
	}
}
