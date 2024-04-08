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

package filehelper

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ******** Public functions ********

// CloseFile closes a file and prints an error message if closing failed.
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		printFileOperationError(`clos`, file.Name(), err)
	}
}

// DeleteFile deletes the file specified by the given file path.
func DeleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		printFileOperationError(`delet`, filePath, err)
	}
}

// GetRealBaseName gets the base name of a file without the extension.
func GetRealBaseName(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}

// FileSize returns the size of the named file.
func FileSize(filePath string) (int64, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

// IsDir checks if a file path is a directory.
func IsDir(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

// IsFileName returns "true", if the supplied file path does not contain any path elements.
func IsFileName(filePath string) bool {
	return filepath.Base(filePath) == filePath
}

// ******** Private functions ********

// printFileOperationError prints an error message for a file operation.
func printFileOperationError(opName string, filePath string, err error) {
	log.Printf(`Error %sing file '%s': %v`, opName, filePath, err)
}
