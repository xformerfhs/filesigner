package numberhelper

// ByteCountForNumber counts how many bytes are needed to represent the given number.
func ByteCountForNumber(number uint64) byte {
	border := uint64(0x00ffffffffffffff)
	byteLen := byte(7)
	for {
		if number > border {
			break
		}
		border >>= 8
		byteLen--
	}

	return byteLen + 1
}
