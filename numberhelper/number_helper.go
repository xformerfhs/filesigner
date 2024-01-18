package numberhelper

func ByteCountForUint64(number uint64) uint {
	border := uint64(0x00ffffffffffffff)
	byteLen := uint(7)
	for {
		if number > border {
			break
		}
		border >>= 8
		byteLen--
	}

	return byteLen + 1
}

func ByteCountForUint32(number uint32) uint {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForUint(number uint) uint {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForInt64(number int64) uint {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForInt32(number int32) uint {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForInt(number int) uint {
	return ByteCountForUint64(uint64(number))
}
