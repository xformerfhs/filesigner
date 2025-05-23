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
// Version: 0.91.0
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

// ******** Private constants ********

// myName contains the program name.
var myName string

// myVersion contains the program version.
const myVersion = `0.91.0`

// mainMsgBase is the base number for all messages in main.
// This file reserves numbers 10-19.
const mainMsgBase = 10

// ******** Formal main function ********

// main is the main function and only a stub for a real main function.
func main() {
	myName = filehelper.GetRealBaseName(os.Args[0])
	// Hack, so that we have a way to have args as arguments, set the exit code and run defer functions.
	// This is a severe design deficiency of Go 1
	os.Exit(mainWithReturnCode(os.Args[1:]))
}

// ******** Private constants ********

// -------- Return codes --------

const (
	rcOK               = 0
	rcCommandLineError = 1
	rcProcessWarning   = 2
	rcProcessError     = 3
)

// -------- Command verbs --------

const (
	commandHelp   = `help`
	commandSign   = `sign`
	commandVerify = `verify`
)

// ******** More private variables ********

// Program information

// scl contains the command line interpreter for the "sign" command.
var scl = cmdline.NewSignCommandLine()

// vcl contains the command line interpreter for the "verify" command.
var vcl = cmdline.NewVerifyCommandLine()

// ******** Real main function ********

// mainWithReturnCode is the real main function with arguments and return code.
// args do not include the program name, only the arguments.
func mainWithReturnCode(args []string) int {
	printVersion()

	argLen := len(args)
	if argLen < 1 {
		return printNotEnoughArgumentsError()
	}

	command := strings.ToLower(args[0])

	switch command {
	case commandHelp:
		printUsageText()
		return rcOK

	case commandSign:
		if argLen < 2 {
			return printMissingArgument(`Context id`)
		}

		contextId := args[1]
		if len(contextId) == 0 {
			printEmptyArgument(`Context id`)
			return rcCommandLineError
		}

		rc := processCmdLineArguments(scl, args[2:])
		if rc != rcOK {
			return rc
		}

		if len(scl.FileList) == 0 {
			logger.PrintWarning(mainMsgBase+2, `No files found to sign`)
			return rcProcessWarning
		}

		return doSigning(scl.SignaturesFileName, scl.SignatureType, contextId, scl.FileList)

	case commandVerify:
		if argLen < 2 {
			return printMissingArgument(`Verification id`)
		}

		verificationId := strings.TrimSpace(args[1])
		if len(verificationId) == 0 {
			printEmptyArgument(`Verification id`)
			return rcCommandLineError
		}

		rc := processCmdLineArguments(vcl, args[2:])
		if rc != rcOK {
			return rc
		}

		return doVerification(vcl.SignaturesFileName, verificationId)

	default:
		return printUsageErrorf(mainMsgBase+3, `Unknown command: '%s'`, command)
	}
}

// ******** Private functions ********

// processCmdLineArguments processes a cmdline.CommandLiner.
func processCmdLineArguments(cl cmdline.CommandLiner, args []string) int {
	err, isHelp := cl.Parse(args)
	if isHelp {
		return rcOK
	}
	if err != nil {
		return printCommandLineParsingError(err)
	}

	err = cl.ExtractCommandData()
	if err != nil {
		logger.PrintErrorf(mainMsgBase+4, `Error getting data from command line: %v`, err)
		return rcProcessError
	}

	return rcOK
}

// printVersion prints the program version information.
func printVersion() {
	logger.PrintInfof(mainMsgBase+5, `%s V%s (%s, %d cpus)`,
		myName,
		myVersion,
		runtime.Version(),
		runtime.NumCPU())
}

// printNotEnoughArgumentsError prints an error message that there are not enough arguments.
func printNotEnoughArgumentsError() int {
	return printUsageError(mainMsgBase+6, `Not enough arguments`)
}

// printMissingArgument prints an error message that the argument with the given name is missing.
func printMissingArgument(name string) int {
	return printUsageErrorf(mainMsgBase+7, `%s is missing`, name)
}

// printEmptyArgument prints an error message that the argument with the given name must not be empty.
func printEmptyArgument(name string) int {
	return printUsageErrorf(mainMsgBase+7, `%s must not be empty`, name)
}

// printCommandLineParsingError prints an error message when there was an error in the command line parsing.
func printCommandLineParsingError(err error) int {
	return printUsageErrorf(mainMsgBase+8, `Error parsing command line: %v`, err)
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
	_, _ = fmt.Printf(`  %s verify {verificationId} [flag]`, myName)
	_, _ = fmt.Print(`

  with 'flag' being the following:

`)
	vcl.PrintUsage()
	_, _ = fmt.Print(`
  The 'verificationId' is the verification id printed when the signatures were created.
  All the files in the signatures file will be verified.


Help:
`)
	_, _ = fmt.Printf(`  %s help`, myName)
	_, _ = fmt.Print(`

  Print this usage information.
`)
}
