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

// SideBand decodes a sideband message from the given data.
func SideBand(data []byte) (SideBandCode, []byte) {
	if len(data) < 1 {
		return SideBandInvalid, data
	}
	code := SideBandCode(data[0])
	switch code {
	case SideBandPackData, SideBandProgress, SideBandFatal:
		return code, data[1:]
	default:
		return SideBandInvalid, data
	}
}

// AppendSideBand appends a sideband message to the given data.
func AppendSideBand(code SideBandCode, data []byte) []byte {
	return append([]byte{byte(code)}, data...)
}
