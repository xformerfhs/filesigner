package numberhelper

func ByteCountForUint64(number uint64) byte {
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

func ByteCountForUint32(number uint32) byte {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForUint(number uint) byte {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForInt64(number int64) byte {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForInt32(number int32) byte {
	return ByteCountForUint64(uint64(number))
}

func ByteCountForInt(number int) byte {
	return ByteCountForUint64(uint64(number))
}
