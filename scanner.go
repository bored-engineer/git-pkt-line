package pktline

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// error-line = PKT-LINE("ERR" SP explanation-text)
// https://git-scm.com/docs/pack-protocol#_pkt_line_format
type ErrErrorLine struct {
	Explanation string
}

func (err ErrErrorLine) Error() string { return "ERR " + err.Explanation }

// pkt-len = 4*(HEXDIG)
type ErrInvalidLen struct {
	Len [4]byte
}

func (err ErrInvalidLen) Error() string {
	return fmt.Sprintf("invalid length prefix: %q", string(err.Len[:]))
}

// 0000: Flush Packet (flush-pkt) - indicates the end of a message
// https://git-scm.com/docs/gitprotocol-v2#_packet_line_framing
var ErrFlushPkt = errors.New("flush-pkt")

// 0001: Delimiter Packet (delim-pkt) - separates sections of a message
// https://git-scm.com/docs/gitprotocol-v2#_packet_line_framing
var ErrDelimPkt = errors.New("delim-pkt")

// 0002: Response End Packet (response-end-pkt) - indicates the end of a response for stateless connections
// https://git-scm.com/docs/gitprotocol-v2#_packet_line_framing
var ErrResponseEndPkt = errors.New("response-end-pkt")

// Scanner parses pkt-line formatted data from an io.Reader
// https://git-scm.com/docs/protocol-common#_pkt_line_format
type Scanner struct {
	r io.Reader
	// Per the documentation the maximum length of a pkt-line’s data component is 65516 bytes
	// However, some implementations can incorrectly send pkt-lines with a length of 65520 bytes
	// So we arrive at 65524 bytes as the maximum buffer size when you include the 4 length bytes
	buf [65524]byte
}

// Reset resets the scanner to read from the given io.Reader
func (s *Scanner) Reset(r io.Reader) { s.r = r }

// Scan advances the Scanner to the next pkt-line
func (s *Scanner) Scan() (line []byte, err error) {
	// Read the length of the next pkt-line
	if n, err := io.ReadFull(s.r, s.buf[0:4]); err != nil {
		if n == 0 && err == io.ErrUnexpectedEOF {
			return nil, io.EOF
		}
		return nil, err
	}
	// Parse the length of the pkt-line as a uint16
	b0, ok0 := unhex(s.buf[0])
	b1, ok1 := unhex(s.buf[1])
	b2, ok2 := unhex(s.buf[2])
	b3, ok3 := unhex(s.buf[3])
	if !ok0 || !ok1 || !ok2 || !ok3 {
		return nil, ErrInvalidLen{
			Len: [4]byte{s.buf[0], s.buf[1], s.buf[2], s.buf[3]},
		}
	}
	len := uint16(b3) | uint16(b2)<<4 | uint16(b1)<<8 | uint16(b0)<<12
	// Verify the legnth is valid (or is a known special packet)
	if len <= 4 || len > 65524 {
		switch len {
		case 0:
			return nil, ErrFlushPkt // 0000
		case 1:
			return nil, ErrDelimPkt // 0001
		case 2:
			return nil, ErrResponseEndPkt // 0002
		default:
			return nil, ErrInvalidLen{
				Len: [4]byte{s.buf[0], s.buf[1], s.buf[2], s.buf[3]},
			}
		}
	}
	// Read the payload of the pkt-line
	if _, err := io.ReadFull(s.r, s.buf[4:len]); err != nil {
		return nil, err
	}
	// Special case, if we see "ERR " at the start of the payload it's an error line
	if explanation, ok := bytes.CutPrefix(s.buf[4:len], []byte("ERR ")); ok {
		return nil, ErrErrorLine{
			Explanation: string(explanation),
		}
	}
	// Success!
	return s.buf[4:len], nil
}

// NewScanner returns a Scanner for a given io.Reader
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: r}
}
