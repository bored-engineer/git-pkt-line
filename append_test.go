package pktline

import (
	"bytes"
	"testing"
)

func TestAppendBytes(t *testing.T) {
	buf := []byte("existing data")
	buf = AppendBytes(buf, []byte("version 1\n"))
	if !bytes.Equal(buf, []byte("existing data000eversion 1\n")) {
		t.Fatalf("AppendBytes did not append the data correctly: %q", string(buf))
	}
}
