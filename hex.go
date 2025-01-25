package pktline

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

// hex converts a byte to its ASCII hex representation
func hex(n byte) byte {
	if n < 10 {
		return '0' + n
	}

	return 'a' - 10 + n
}
