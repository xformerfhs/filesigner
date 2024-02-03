package maphelper

import (
	"golang.org/x/exp/maps"
	"golang.org/x/text/cases"
)

var foldCaser = cases.Fold()

type CaseInsensitiveStringMap[K any] struct {
	m map[string]K
}

func NewCaseInsensitiveStringMap[K any]() *CaseInsensitiveStringMap[K] {
	return &CaseInsensitiveStringMap[K]{}
}

func (cism *CaseInsensitiveStringMap[K]) Add(k string, v K) {
	foldedKey := foldCaser.String(k)
	cism.m[foldedKey] = v
}

func (cism *CaseInsensitiveStringMap[K]) Remove(k string) {
	foldedKey := foldCaser.String(k)
	delete(cism.m, foldedKey)
}

func (cism *CaseInsensitiveStringMap[K]) Contains(k string) bool {
	foldedKey := foldCaser.String(k)
	_, exists := cism.m[foldedKey]
	return exists
}

func (cism *CaseInsensitiveStringMap[K]) Len() int {
	return len(cism.m)
}

func (cism *CaseInsensitiveStringMap[K]) Clear() {
	clear(cism.m)
}

func (cism *CaseInsensitiveStringMap[K]) Keys() []string {
	return maps.Keys(cism.m)
}

func (cism *CaseInsensitiveStringMap[K]) Values() []K {
	return maps.Values(cism.m)
}
