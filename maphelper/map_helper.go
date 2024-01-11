package maphelper

import (
	"cmp"
	"golang.org/x/exp/maps"
	"slices"
)

// SortedKeys returns the keys of the map m. The keys will be sorted.
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	result := maps.Keys(m)
	slices.Sort(result)

	return result
}
