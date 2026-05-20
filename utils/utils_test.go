package utils

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestListPrint(t *testing.T) {
	list := []string{"hello", "goodbye"}
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ListPrint(list)

	_ = w.Close()
	os.Stdout = originalStdout

	capturedOutput, _ := io.ReadAll(r)

	expected := "hello goodbye"

	if !strings.Contains(string(capturedOutput), expected) {
		t.Errorf("Expected output: %s, got: %s", expected, capturedOutput)
	}

}
