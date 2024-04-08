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
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// expectedGoFilesCount is the no. of *.go files in this directory. Change if count changes.
const expectedGoFilesCount = 8

func TestEmpty(t *testing.T) {
	fileList, err := SensibleGlob("")
	if err != nil {
		t.Fatalf(`Error with an empty search pattern: %v`, err)
	}
	if len(fileList) != 0 {
		t.Fatalf(`File list is not empty with empty search pattern`)
	}
}

func TestOnlyTrailing(t *testing.T) {
	fileList, err := SensibleGlob(string(filepath.Separator))
	if err != nil {
		t.Fatalf(`Error with an only separator pattern: %v`, err)
	}
	if len(fileList) != 0 {
		t.Fatalf(`File list is not empty with only separator pattern`)
	}
}

func TestOnlyManyTrailing(t *testing.T) {
	pattern := strings.Repeat(string(filepath.Separator), 100)
	fileList, err := SensibleGlob(pattern)
	if err != nil {
		t.Fatalf(`Error with a many separators pattern: %v`, err)
	}
	if len(fileList) != 0 {
		t.Fatalf(`File list is not empty with many separators pattern`)
	}
}

func TestFilesManyTrailing(t *testing.T) {
	pattern := `*.go` + strings.Repeat(string(filepath.Separator), 100)
	fileList, err := SensibleGlob(pattern)
	if err != nil {
		t.Fatalf(`Error with a correct pattern with many separators pattern: %v`, err)
	}
	if len(fileList) != expectedGoFilesCount {
		t.Fatalf(`File list has wrong size with many separators`)
	}
}

func TestFilesNormal(t *testing.T) {
	fileList, err := SensibleGlob("*.go")
	if err != nil {
		t.Fatalf(`Error with correct pattern: %v`, err)
	}
	if len(fileList) != expectedGoFilesCount {
		t.Fatalf(`File list has wrong size with correct pattern`)
	}
}

func TestFilesWrongCase(t *testing.T) {
	fileList, err := SensibleGlob("*.GO")
	if err != nil {
		t.Fatalf(`Error with wrong case pattern: %v`, err)
	}

	var expectedCount int
	if runtime.GOOS == `windows` {
		expectedCount = expectedGoFilesCount
	} else {
		expectedCount = 0
	}

	if len(fileList) != expectedCount {
		t.Fatalf(`File list has wrong size with correct pattern`)
	}
}

func TestOneFileNormal(t *testing.T) {
	fileList, err := SensibleGlob("glob_test.go")
	if err != nil {
		t.Fatalf(`Error with one file pattern: %v`, err)
	}

	if len(fileList) != 1 {
		t.Fatalf(`File list has wrong size with one fie pattern`)
	}

	if fileList[0] != `glob_test.go` {
		t.Fatalf(`Wrong file found: %s`, fileList[0])
	}
}

func TestInvalidFilePattern(t *testing.T) {
	fileList, err := SensibleGlob(`*.<|>go`)
	if runtime.GOOS == `windows` && err == nil {
		t.Fatalf(`No error with an invalid pattern`)
	}

	if len(fileList) != 0 {
		t.Fatalf(`Files present with an invalid pattern`)
	}
}

func TestUnknownFilePattern(t *testing.T) {
	fileList, err := SensibleGlob(`*.xyz`)
	if err != nil {
		t.Fatalf(`Error with an unknown pattern: %v`, err)
	}

	if len(fileList) != 0 {
		t.Fatalf(`Files present with an unknown pattern`)
	}
}

func TestFilesOnly(t *testing.T) {
	fileList, err := SensibleGlob("*")
	if err != nil {
		t.Fatalf(`Error with everything pattern: %v`, err)
	}

	for _, fileName := range fileList {
		var fi os.FileInfo
		fi, err = os.Stat(fileName)
		if err != nil {
			t.Fatalf(`Error trying to stat '%s': %v`, fileName, err)
		}

		if fi.IsDir() {
			t.Fatalf(`Glob returned directory: '%s''`, fileName)
		}
	}
}
