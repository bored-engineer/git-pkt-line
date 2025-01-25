package pktline

import (
	"fmt"
)

// 0000: Flush Packet (flush-pkt) - indicates the end of a message
func AppendFlushPkt(b []byte) []byte {
	return append(b, '0', '0', '0', '0')
}

// 0001: Delimiter Packet (delim-pkt) - separates sections of a message
func AppendDelimPkt(b []byte) []byte {
	return append(b, '0', '0', '0', '1')
}

// 0002: Response End Packet (response-end-pkt) - indicates the end of a response for stateless connections
func AppendResponseEndPkt(b []byte) []byte {
	return append(b, '0', '0', '0', '2')
}

// AppendLength appends the given lengthx to the slice as ASCII hex characters
func AppendLength(b []byte, sz int) []byte {
	sz += 4
	if sz > 65520 {
		panic(fmt.Errorf("length (%d) overflows maximum permitted value (65520)", sz))
	}
	return append(b,
		hex(byte(sz&0xf000>>12)),
		hex(byte(sz&0x0f00>>8)),
		hex(byte(sz&0x00f0>>4)),
		hex(byte(sz&0x000f)),
	)
}

// AppendBytes appends the given bytes to the slice, prepending the pkt-len
func AppendBytes(b []byte, data []byte) []byte {
	return append(AppendLength(b, len(data)), data...)
}

// AppendBytes appends the given string to the slice, prepending the pkt-len
func AppendString(b []byte, data string) []byte {
	return append(AppendLength(b, len(data)), data...)
}
