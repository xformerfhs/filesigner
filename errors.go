//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
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
//    2025-05-25: V1.0.0: Created.
//

package main

import (
	"filesigner/logger"
	"fmt"
)

// ******** Private functions ********

// printNotEnoughArgumentsError prints an error message that there are not enough arguments.
func printNotEnoughArgumentsError() int {
	return printUsageError(errorMsgBase+0, `Not enough arguments`)
}

// printMissingArgument prints an error message that the argument with the given name is missing.
func printMissingArgument(name string) int {
	return printUsageErrorf(errorMsgBase+1, `%s is missing`, name)
}

// printEmptyArgument prints an error message that the argument with the given name must not be empty.
func printEmptyArgument(name string) int {
	return printUsageErrorf(errorMsgBase+2, `%s must not be empty`, name)
}

// printCommandLineParsingError prints an error message when there was an error in the command line parsing.
func printCommandLineParsingError(err error) int {
	return printUsageErrorf(errorMsgBase+3, `Error parsing command line: %v`, err)
}

// printUsageError prints an error message followed by the usage message.
func printUsageError(msgNum byte, msgText string) int {
	logger.PrintError(msgNum, msgText)
	printUsageText()
	return rcCommandLineError
}

// printUsageErrorf prints an error message followed by the usage message with a format string.
func printUsageErrorf(msgNum byte, msgFormat string, args ...any) int {
	logger.PrintErrorf(msgNum, msgFormat, args...)
	printUsageText()
	return rcCommandLineError
}

// printUsageText prints the usage text.
func printUsageText() {
	_, _ = fmt.Print(`
Usage:
  Create and verify signatures for a collection of files.


Sign files:
`)
	_, _ = fmt.Printf(`  %s sign {contextId} [flags] [files]`, myName)
	_, _ = fmt.Print(`

  with 'files' being an optional list of file names and 'flags' one or more of the following options:

`)
	scl.PrintUsage()
	_, _ = fmt.Print(`
  The 'contextId' is an arbitrary word used to make the signature depend on a topic, also called a 'domain separator'.
  If no file names are specified, all files in the current directory are signed.
  This can be modified by the exclude and include options.
  The '--recurse' option is only valid if there are either no files specified or if there are include options present.
  The files must be present in the current directory or one of its subdirectories.
  Specifying a file outside the current directory tree is an error.
  All file names that contain wildcards ('*', '?') are treated as if they were specified in an '--include-file' option.


Verify files:
`)
	_, _ = fmt.Printf(`  %s verify {verificationId} [flags]`, myName)
	_, _ = fmt.Print(`

  with 'flags' being one or more of the following options:

`)
	vcl.PrintUsage()
	_, _ = fmt.Print(`
  The 'verificationId' is the verification id printed when the signatures were created.
  All the files in the signatures file will be verified.


Get version:
`)
	_, _ = fmt.Printf(`  %s version`, myName)
	_, _ = fmt.Print(`

  Print version information.


Help:
`)
	_, _ = fmt.Printf(`  %s help`, myName)
	_, _ = fmt.Print(`

  Print this usage information.
`)
}
