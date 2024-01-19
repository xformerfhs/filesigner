package numberhelper

import (
	"errors"
	"filesigner/slicehelper"
)

// ******** Private constants ********

const maxByteCounterLen = 8

// ******** Public types ********

type ByteSliceCounter struct {
	counter []byte
	first   byte
	last    byte
}

// ******** Public creation functions ********

func NewByteSliceCounter(length byte) (*ByteSliceCounter, error) {
	if length == 0 {
		return nil, errors.New("Byte counter length must not be 0")
	}

	if length > maxByteCounterLen {
		return nil, errors.New("Byte counter length is too large")
	}

	return &ByteSliceCounter{counter: make([]byte, length), first: length - 1, last: length - 1}, nil
}

func NewByteSliceCounterForCount(count uint) (*ByteSliceCounter, error) {
	return NewByteSliceCounter(ByteCountForNumber(uint64(count)))
}

// ******** Public functions ********

// -------- Counting functions --------

func (bc *ByteSliceCounter) Inc() {
	counter := bc.counter
	i := bc.last
	for {
		a := counter[i]
		a++
		counter[i] = a
		if a != 0 || i == 0 {
			break
		}

		i--
	}

	bc.first = i
}

func (bc *ByteSliceCounter) Dec() {
	counter := bc.counter
	i := byte(len(counter) - 1)
	for {
		a := counter[i]
		a--
		counter[i] = a
		if a != 0xff || i == 0 {
			break
		}

		i--
	}
}

func (bc *ByteSliceCounter) Clear() {
	counter := bc.counter
	slicehelper.ClearInteger(counter)
}

func (bc *ByteSliceCounter) Slice() []byte {
	return bc.counter[bc.first:]
}

// -------- Setter and getter functions --------

func (bc *ByteSliceCounter) SetUint64(value uint64) {
	counter := bc.counter
	i := bc.last
	for {
		counter[i] = byte(value)
		value >>= 8
		if value == 0 {
			break
		}

		i--
	}

	bc.first = i
}

func (bc *ByteSliceCounter) GetUint64() uint64 {
	result := uint64(0)
	first := bc.first
	counter := bc.counter
	i := bc.last
	for {
		result <<= 8
		result |= uint64(counter[i])

		if i == first {
			break
		}
	}

	return result
}
