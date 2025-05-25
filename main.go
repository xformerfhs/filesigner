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
// Version: 0.92.0
//

package main

import (
	"filesigner/cmdline"
	"filesigner/filehelper"
	"filesigner/logger"
	"os"
	"runtime"
	"strings"
)

// ******** Private constants ********

// myName contains the program name.
var myName string

// myVersion contains the program version.
const myVersion = `0.92.0`

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
	commandHelp    = `help`
	commandSign    = `sign`
	commandVerify  = `verify`
	commandVersion = `version`
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
		if len(args) < 2 {
			return printMissingArgument(`Context id`)
		}
		return handleSign(args[1:])

	case commandVerify:
		if len(args) < 2 {
			return printMissingArgument(`Verification id`)
		}
		return handleVerify(args[1:])

	case commandVersion:
		return printVersion()

	default:
		return printUsageErrorf(mainMsgBase+0, `Unknown command: '%s'`, command)
	}
}

// ******** Private functions ********

// printVersion prints the program version information.
func printVersion() int {
	logger.PrintInfof(mainMsgBase+1, `%s V%s (%s, %d cpus)`,
		myName,
		myVersion,
		runtime.Version(),
		runtime.NumCPU())

	return rcOK
}
