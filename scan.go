package pktline

import (
	"bytes"
	"fmt"
	"io"
	"iter"
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
type ErrFlushPkt struct{}

func (err ErrFlushPkt) Error() string { return "flush-pkt" }

// 0001: Delimiter Packet (delim-pkt) - separates sections of a message
// https://git-scm.com/docs/gitprotocol-v2#_packet_line_framing
type ErrDelimPkt struct{}

func (err ErrDelimPkt) Error() string { return "delim-pkt" }

// 0002: Response End Packet (response-end-pkt) - indicates the end of a response for stateless connections
// https://git-scm.com/docs/gitprotocol-v2#_packet_line_framing
type ErrResponseEndPkt struct{}

func (err ErrResponseEndPkt) Error() string { return "delim-pkt" }

// unhex converts an ascii hex character to its numeric value
func unhex(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// Parses a pkt-line length without using the strconv or hex package
func parsePktLen(len [4]byte) (uint16, error) {
	b0, ok0 := unhex(len[0])
	b1, ok1 := unhex(len[1])
	b2, ok2 := unhex(len[2])
	b3, ok3 := unhex(len[3])
	if !ok0 || !ok1 || !ok2 || !ok3 {
		return 0, ErrInvalidLen{len}
	}
	sz := uint16(b3) | uint16(b2)<<4 | uint16(b1)<<8 | uint16(b0)<<12
	switch sz {
	case 0:
		return 0, ErrFlushPkt{}
	case 1:
		return 0, ErrDelimPkt{}
	case 2:
		return 0, ErrResponseEndPkt{}
	case 3:
		return 0, ErrInvalidLen{len}
	default:
		return sz - 4, nil
	}
}

// Scan reads pkt-line packets from the reader until an error occurs
func Scan(r io.Reader) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		var pktLenBytes [4]byte
		var payload []byte
		for {
			// Read and parse the length of the next packet
			if _, err := io.ReadFull(r, pktLenBytes[:]); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				if !yield(nil, err) {
					return
				}
				continue
			}
			pktLen, err := parsePktLen(pktLenBytes)
			if err != nil {
				if !yield(nil, err) {
					return
				}
				continue
			}

			// If we can't fit the payload in the capacity of the buffer, expand it
			if cap(payload) < int(pktLen) {
				payload = append(payload[:cap(payload)], make([]byte, int(pktLen)-cap(payload))...)
			}

			// Set the slice length to the size of the payload
			payload = payload[:pktLen]

			// Read the full packet/payload
			if _, err := io.ReadFull(r, payload); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				if !yield(nil, err) {
					return
				}
				continue
			}

			// Special case, if we see "ERR " at the start of the payload it's an error line
			if explanation, ok := bytes.CutPrefix(payload, []byte("ERR ")); ok {
				if !yield(nil, ErrErrorLine{
					Explanation: string(explanation),
				}) {
					return
				}
				continue
			}

			// Otherwise it's a regular packet
			if !yield(payload, nil) {
				return
			}

		}
	}
}
