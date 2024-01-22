package typelist

import (
	"fmt"
)

type FlagTypeList []string

func NewFlagTypeList() *FlagTypeList {
	result := make(FlagTypeList, 0, 100)
	return &result
}

func (ftl *FlagTypeList) String() string {
	return fmt.Sprint(*ftl)
}

func (ftl *FlagTypeList) Set(value string) error {
	*ftl = append(*ftl, value)

	return nil
}

func (ftl *FlagTypeList) GetNames() []string {
	return *ftl
}

func (ftl *FlagTypeList) Len() int {
	return len(*ftl)
}

func (ftl *FlagTypeList) HasEntries() bool {
	return len(*ftl) > 0
}
