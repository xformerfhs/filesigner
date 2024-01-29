package slicehelper

import "golang.org/x/exp/constraints"

// ******** Public functions ********

// Fill fills a generic slice with a generic value in an efficient way.
func Fill[S ~[]T, T any](a S, value T) {
	aLen := ensureLengthIsCapacity(&a)

	if aLen > 0 {
		// Put the value into the first slice element
		a[0] = value

		// Incrementally duplicate the value into the rest of the slice
		for j := 1; j < aLen; j <<= 1 {
			copy(a[j:], a[:j])
		}
	}
}

// ClearInteger clears an integer type slice.
func ClearInteger[S ~[]T, T constraints.Integer](a S) {
	Fill(a, 0)
}

// MakeCopy makes a copy of slice
func MakeCopy[S ~[]T, T any](a S) S {
	result := make([]T, len(a))
	copy(result, a)
	return result
}

// Prepend adds an element v at the beginning of a slice s.
func Prepend[T any](v T, s []T) []T {
	return append([]T{v}, s...)
}

// ******** Private functions ********

// ensureLengthIsCapacity ensures that the length of the slice is its capacity.
// We need the address of the slice as the parameter. If the '*' would be missing
// we would get a copy of the slice and not the slice itself.
func ensureLengthIsCapacity[S ~[]T, T any](a *S) int {
	ra := *a
	aLen := len(ra)
	aCap := cap(ra)
	if aLen != aCap {
		*a = ra[:aCap]
		aLen = aCap
	}

	return aLen
}
