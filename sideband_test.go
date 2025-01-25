package pktline

import "testing"

func TestSideband(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    SideBandCode
		wantBuf string
	}{
		"empty": {
			input:   "",
			want:    SideBandInvalid,
			wantBuf: "",
		},
		"invalid": {
			input:   "\xffinvalid",
			want:    SideBandInvalid,
			wantBuf: "\xffinvalid",
		},
		"pack": {
			input:   "\x01PACK",
			want:    SideBandPackData,
			wantBuf: "PACK",
		},
		"progress": {
			input:   "\x02progress",
			want:    SideBandPackData,
			wantBuf: "progress",
		},
		"fatal": {
			input:   "\x03fatal",
			want:    SideBandFatal,
			wantBuf: "fatal",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, buf := SideBand([]byte(tc.input))
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
			if string(buf) != tc.wantBuf {
				t.Errorf("got %v, want %v", string(buf), tc.wantBuf)
			}
		})
	}
}
