package pktline

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := map[string]struct {
		input    string
		want     []byte
		wantErr  error
		wantScan bool
	}{
		"eof": {
			input:   "",
			wantErr: io.EOF,
		},
		"unexpected-eof-length": {
			input:   "000",
			wantErr: io.ErrUnexpectedEOF,
		},
		"unexpected-eof-payload": {
			input:   "001c incomplete",
			wantErr: io.ErrUnexpectedEOF,
		},
		"invalid-hex-length": {
			input: "$$$$",
			wantErr: ErrInvalidLen{
				Len: [4]byte{'$', '$', '$', '$'},
			},
		},
		"flush-pkt": {
			input:   "0000",
			wantErr: ErrFlushPkt,
		},
		"delim-pkt": {
			input:   "0001",
			wantErr: ErrDelimPkt,
		},
		"response-end-pkt": {
			input:   "0002",
			wantErr: ErrResponseEndPkt,
		},
		"invalid-length": {
			input: "0003",
			wantErr: ErrInvalidLen{
				Len: [4]byte{'0', '0', '0', '3'},
			},
		},
		"err-line": {
			input: "001cERR something went wrong",
			wantErr: ErrErrorLine{
				Explanation: "something went wrong",
			},
		},
		"valid": {
			input: "000eversion 1\n",
			want:  []byte("version 1\n"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var scanner Scanner
			scanner.Reset(strings.NewReader(tc.input))
			gotPayload, gotErr := scanner.Scan()
			if tc.wantErr == nil && gotErr != nil {
				t.Fatalf("expected no error, got %v", gotErr)
			} else if gotErr != tc.wantErr {
				t.Fatalf("expected error %v, got %v", tc.wantErr, gotErr)
			}
			if !bytes.Equal(gotPayload, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, gotPayload)
			}
		})
	}
}
