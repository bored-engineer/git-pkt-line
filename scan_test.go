package pktline

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    []byte
		wantErr error
	}{
		"unexpected-eof": {
			input:   "000",
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
			wantErr: ErrFlushPkt{},
		},
		"delim-pkt": {
			input:   "0001",
			wantErr: ErrDelimPkt{},
		},
		"response-end-pkt": {
			input:   "0002",
			wantErr: ErrResponseEndPkt{},
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
		"empty": {
			input: "0004",
			want:  []byte{},
		},
		"valid": {
			input: "000eversion 1\n",
			want:  []byte("version 1\n"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for got, err := range Scan(strings.NewReader(tc.input)) {
				if tc.wantErr != nil {
					if err == nil {
						t.Fatalf("expected error, got nil")
					} else if !errors.Is(err, tc.wantErr) {
						t.Fatalf("expected error %v, got %v", tc.wantErr, err)
					}
				} else {
					if err != nil {
						t.Fatalf("expected error to be nil, got %v", err)
					}
					if !bytes.Equal(got, tc.want) {
						t.Fatalf("expected %v, got %v", tc.want, got)
					}
				}
				break
			}
		})
	}
}
