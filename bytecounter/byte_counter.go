package bytecounter

import (
	"errors"
	"filesigner/numberhelper"
	"filesigner/slicehelper"
)

const maxByteCounterLen = 16

type ByteSliceCounter struct {
	Counter []byte
}

func NewByteSliceCounter(length byte) (*ByteSliceCounter, error) {
	if length == 0 {
		return nil, errors.New("Byte counter length must not be 0")
	}

	if length > maxByteCounterLen {
		return nil, errors.New("Byte counter length is too large")
	}

	return &ByteSliceCounter{Counter: make([]byte, length)}, nil
}

func NewByteSliceCounterForCount(count uint) (*ByteSliceCounter, error) {
	return NewByteSliceCounter(numberhelper.ByteCountForUint(count))
}

func (bc *ByteSliceCounter) Inc() {
	counter := bc.Counter
	for i := len(counter) - 1; i >= 0; i-- {
		a := counter[i]
		a++
		counter[i] = a
		if a != 0 {
			break
		}
	}
}

func (bc *ByteSliceCounter) Dec() {
	counter := bc.Counter
	for i := len(counter) - 1; i >= 0; i-- {
		a := counter[i]
		a--
		counter[i] = a
		if a != 0xff {
			break
		}
	}
}

func (bc *ByteSliceCounter) Clear() {
	counter := bc.Counter
	slicehelper.ClearInteger(counter)
}
