//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-02-05: V2.0.0: Rewritten and renamed to use FileSystemStringSet.
//

package flaglist

import (
	"filesigner/set"
	"fmt"
)

// FileSystemFlagList collects names from flags that are supposed to be file system names.
// It implements the "Value" interface.
type FileSystemFlagList struct {
	st *set.FileSystemStringSet
}

// NewFileSystemFlagList returns an empty FileSystemFlagList
func NewFileSystemFlagList() *FileSystemFlagList {
	return &FileSystemFlagList{st: set.NewFileSystemStringSet()}
}

// String returns the string representation of the elements of the list.
// It is part of the "Value" interface.
func (fl *FileSystemFlagList) String() string {
	return fmt.Sprint(fl.st.Elements())
}

// Set adds a value to the list.
// It is part of the "Value" interface.
func (fl *FileSystemFlagList) Set(value string) error {
	fl.st.Add(value)

	return nil
}

// Elements returns the elements of the list.
func (fl *FileSystemFlagList) Elements() []string {
	return fl.st.Elements()
}

// Len returns the number of elements in the list.
func (fl *FileSystemFlagList) Len() int {
	return fl.st.Len()
}

// HasElements returns "true", if there are elements in the list, "false" otherwise.
func (fl *FileSystemFlagList) HasElements() bool {
	return fl.st.Len() != 0
}
