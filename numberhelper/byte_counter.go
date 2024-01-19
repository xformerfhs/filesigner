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

	length--

	return &ByteSliceCounter{counter: make([]byte, length), first: length, last: length}, nil
}

func NewByteSliceCounterForCount(count uint64) (*ByteSliceCounter, error) {
	return NewByteSliceCounter(ByteCountForNumber(count))
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

func (bc *ByteSliceCounter) Zero() {
	counter := bc.counter
	slicehelper.ClearInteger(counter)
}

func (bc *ByteSliceCounter) Slice() []byte {
	return bc.counter[bc.first:]
}

func (bc *ByteSliceCounter) FullSlice() []byte {
	return bc.counter
}

// -------- Setter and getter functions --------

func (bc *ByteSliceCounter) SetCount(value uint64) {
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

func (bc *ByteSliceCounter) GetCount() uint64 {
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

		i--
	}

	return result
}
