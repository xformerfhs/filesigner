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
//

package filehasher

import (
	"filesigner/filehelper"
	"filesigner/numberhelper"
	"filesigner/paddedhasher"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
)

// ******** Public types ********

type FileHasher struct {
	hasher *paddedhasher.PaddedHasher
}

// ******** Creation functions ********

// NewFileHasher Create new file hasher structure.
func NewFileHasher(contextKey []byte) (*FileHasher, error) {
	hasher := paddedhasher.NewPaddedHasher(sha3.New512(), contextKey)

	return &FileHasher{hasher}, nil
}

// ******** Public functions ********

// HashFile calculates the hash value for one file.
func (fh *FileHasher) HashFile(filePath string) ([]byte, error) {
	hasher := fh.hasher

	err := hashFileContent(hasher, filePath)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

// ******** Private functions ********

// hashFileContent writes the content of a file to a hasher.
func hashFileContent(hasher *paddedhasher.PaddedHasher, filePath string) error {
	var err error
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return err
	}
	defer filehelper.CloseFile(f)

	_, err = io.Copy(hasher, f)
	if err != nil {
		return err
	}

	_, err = hasher.Write(numberhelper.Uint64AsShortestBigEndianBytes(hasher.Count()))
	if err != nil {
		return err
	}

	return nil
}
