package maphelper

import (
	"cmp"
	"golang.org/x/exp/maps"
	"slices"
)

// SortedKeys returns a slice with the keys of a map ordered by value
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	result := maps.Keys(m)
	slices.Sort(result)

	return result
}
