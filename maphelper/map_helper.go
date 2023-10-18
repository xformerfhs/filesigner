package maphelper

import (
	"golang.org/x/exp/maps"
	"sort"
)

// GetSortedKeys sorts the keys of a map with string keys
func GetSortedKeys[T any](a map[string]T) []string {
	result := maps.Keys(a)
	sort.Strings(result)

	return result
}
