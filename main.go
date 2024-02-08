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

package main

import (
	"filesigner/cmdline"
	"filesigner/filehelper"
	"filesigner/logger"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// ******** Formal main function ********

// main is the main function and only a stub for a real main function.
func main() {
	// Hack, so that we have a way to have args as arguments, set the exit code and run defer functions.
	// This is a severe design deficiency of Go 1
	os.Exit(mainWithReturnCode(os.Args))
}

// Private constants

// Return codes

const (
	rcOK               = 0
	rcCommandLineError = 1
	rcProcessWarning   = 2
	rcProcessError     = 3
)

// Command verbs

const (
	commandHelp   = "help"
	commandSign   = "sign"
	commandVerify = "verify"
)

// Program information

// myVersion contains the program version.
const myVersion = "0.70.0"

// myName contains the program name.
var myName string

// scl contains the command line interpreter for the "sign" command.
var scl = cmdline.NewSignCommandLine()

// vcl contains the command line interpreter for the "verify" command.
var vcl = cmdline.NewVerifyCommandLine()

// ******** Real main function ********

// mainWithReturnCode is the real main function with arguments and return code.
func mainWithReturnCode(args []string) int {
	myName = filehelper.GetRealBaseName(args[0])

	printVersion()

	argLen := len(args)
	if argLen < 2 {
		return printUsageError(11, `not enough arguments`)
	}

	command := strings.ToLower(args[1])

	switch command {
	case commandHelp:
		printUsageText()
		return rcOK

	case commandSign:
		if argLen < 3 {
			return printMissingContextId()
		}

		err, isHelp := scl.Parse(args[3:])
		if isHelp {
			return rcOK
		}
		if err != nil {
			return printCommandLineProcessingError(err)
		}

		err = scl.ExtractCommandData()
		if err != nil {
			logger.PrintError(12, err.Error())
			return rcProcessError
		}

		return doSigning(scl.SignaturesFileName, scl.SignatureType, args[2], scl.FileList)

	case commandVerify:
		if argLen < 3 {
			return printMissingContextId()
		}

		err, isHelp := vcl.Parse(args[3:])
		if isHelp {
			return rcOK
		}
		if err != nil {
			return printCommandLineProcessingError(err)
		}

		return doVerification(args[2], vcl.SignaturesFileName)

	default:
		return printUsageErrorf(13, `unknown command: '%s'`, command)
	}
}

// ******** Private functions ********

// printVersion prints the program version information.
func printVersion() {
	logger.PrintInfof(14, "%s V%s (%s, %d cpus)",
		myName,
		myVersion,
		runtime.Version(),
		runtime.NumCPU())
}

func printMissingContextId() int {
	return printUsageError(15, "context id missing")
}

// printCommandLineProcessingError prints an error message when there was in error in the command line processing.
func printCommandLineProcessingError(err error) int {
	return printUsageErrorf(16, "error processing command line: %v", err)
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
	_, _ = fmt.Println("\nUsage:")

	_, _ = fmt.Println("\nSign files:")
	_, _ = fmt.Printf("  %s sign {contextId} [flags] [files]\n", myName)
	_, _ = fmt.Println("\n  with 'flags' being one of the following and 'files' a list of file names:\n")
	scl.PrintUsage()
	_, _ = fmt.Println("\n  If no file names are given, the current directory is searched for files.")
	_, _ = fmt.Println(`  This can be modified by the exclude and include options.`)
	_, _ = fmt.Println(`  If no files and no exclude/include options are present, all files in the current directory will be signed.`)
	_, _ = fmt.Println(`  The '--recurse' option is only valid if there are either no files specified or if there are include options present.`)
	_, _ = fmt.Println("  The files must be present in the current directory or one of its subdirectories.")
	_, _ = fmt.Println("  Specifying a file outside the current directory tree is an error.")
	_, _ = fmt.Println("  All file names that contain wildcards ('*', '?') are treated as if they were specified in an '--include-file' option.\n")

	_, _ = fmt.Println("\nVerify files:")
	_, _ = fmt.Printf("  %s verify {contextId} [flag]\n", myName)
	_, _ = fmt.Println("\n  with 'flag' being the following:\n")
	vcl.PrintUsage()
	_, _ = fmt.Println("\n  All the files in the signatures file will be verified.\n")

	_, _ = fmt.Println(`  The 'contextId' is an arbitrary word used to make the signature depend on a topic, also called a 'domain separator'.`)
	_, _ = fmt.Println("\n\nHelp:")
	_, _ = fmt.Printf("  Call: %s help\n", myName)
	_, _ = fmt.Println("\n  Print this usage information.")
}
