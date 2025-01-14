package pktline

// byteToASCIIHex converts a byte to its ASCII hex representation
func byteToASCIIHex(n byte) byte {
	if n < 10 {
		return '0' + n
	}

	return 'a' - 10 + n
}

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

// AppendBytes appends the given bytes to the slice adding a length prefix to the front
func AppendBytes(b []byte, data []byte) []byte {
	sz := len(data) + 4
	return append(
		append(b,
			byteToASCIIHex(byte(sz&0xf000>>12)),
			byteToASCIIHex(byte(sz&0x0f00>>8)),
			byteToASCIIHex(byte(sz&0x00f0>>4)),
			byteToASCIIHex(byte(sz&0x000f)),
		),
		data...,
	)
}

// AppendString is an alias for AppendBytes but takes a string type
func AppendString(b []byte, data string) []byte {
	return AppendBytes(b, []byte(data))
}
