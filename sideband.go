package pktline

type SideBandCode byte

const (
	// Returned when the sideband code is invalid
	SideBandInvalid SideBandCode = 0
	// 1 - pack data
	SideBandPackData SideBandCode = 1
	// 2 - progress messages
	SideBandProgress SideBandCode = 2
	// 3 - fatal error message just before stream aborts
	SideBandFatal SideBandCode = 3
)

// SideBand decodes a sideband message from the given pkt-line.
func SideBand(line []byte) (SideBandCode, []byte) {
	if len(line) < 1 {
		return SideBandInvalid, line
	}
	code := SideBandCode(line[0])
	switch code {
	case SideBandPackData, SideBandProgress, SideBandFatal:
		return code, line[1:]
	default:
		return SideBandInvalid, line
	}
}

// AppendSideBand prepends a sideband code to the given pkt-line.
func AppendSideBand(code SideBandCode, line []byte) []byte {
	return append([]byte{byte(code)}, line...)
}
