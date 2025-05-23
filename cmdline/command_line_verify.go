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
// Version: 2.0.0
//
// Change history:
//    2024-02-08: V1.0.0: Created.
//    2024-04-05: V1.0.1: Make Stdout the output destination for usage messages.
//    2025-05-23: V2.0.0: Add verification id.
//

package cmdline

import (
	"errors"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

// ******** Public types ********

// VerifyCommandLine is the object that contains all the data
// to interpret a "verify" command line.
type VerifyCommandLine struct {
	// Public elements
	SignaturesFileName string
	VerificationId     string

	// Private elements
	fs     *pflag.FlagSet
	prefix string
}

// ******** Public functions ********

// NewVerifyCommandLine sets up the flag parser for the "verify" command.
func NewVerifyCommandLine() *VerifyCommandLine {
	verifyCmd := pflag.NewFlagSet(`verify`, pflag.ContinueOnError)

	verifyCmd.SetOutput(os.Stdout)

	result := &VerifyCommandLine{fs: verifyCmd}

	verifyCmd.StringVarP(&result.prefix, `name`, `m`, defaultSignaturesFileNamePrefix, `Prefix of the signatures file name`)
	verifyCmd.StringVarP(&result.VerificationId, `verificationid`, `v`, ``, `Verification id`)

	verifyCmd.SortFlags = true

	return result
}

// Parse parses the command line according to the flag rules.
func (cl *VerifyCommandLine) Parse(args []string) (error, bool) {
	err := cl.fs.Parse(args)
	if errors.Is(err, pflag.ErrHelp) {
		return nil, true
	}

	if cl.fs.NArg() != 0 {
		return errors.New(`There must be no arguments present without options`), false
	}

	return err, false
}

// PrintUsage prints the usage information for the command.
func (cl *VerifyCommandLine) PrintUsage() {
	cl.fs.PrintDefaults()
}

// ExtractCommandData returns the data that are needed for the command.
func (cl *VerifyCommandLine) ExtractCommandData() error {
	// 1. Build signatures file name.
	cl.SignaturesFileName = cl.prefix + signaturesFileNameSuffix

	// 2. The signatures file must be written to the current directory.
	err := checkSignaturesFileName(cl.SignaturesFileName)
	if err != nil {
		return err
	}

	cl.VerificationId = strings.TrimSpace(cl.VerificationId)
	if len(cl.VerificationId) == 0 {
		return errors.New(`Verification id must not be empty`)
	}

	return nil
}
