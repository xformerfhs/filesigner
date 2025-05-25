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
	"filesigner/cmdline"
	"filesigner/logger"
	"strings"
)

// ******** Private functions ********

// handleSign processes the "sign" command.
func handleSign(args []string) int {
	contextId := args[0]
	if len(contextId) == 0 {
		printEmptyArgument(`Context id`)
		return rcCommandLineError
	}

	rc := processCmdLineArguments(scl, args[1:])
	if rc != rcOK {
		return rc
	}

	if scl.BeQuiet {
		logger.SetLogLevel(logger.LogLevelWarning)
	}

	if len(scl.FileList) == 0 {
		logger.PrintWarning(handlerMsgBase+0, `No files found to sign`)
		return rcProcessWarning
	}

	return doSigning(scl.SignaturesFileName, scl.SignatureType, contextId, scl.BeQuiet, scl.FileList)
}

// handleVerify processes the "verify" command.
func handleVerify(args []string) int {
	verificationId := strings.TrimSpace(args[0])
	if len(verificationId) == 0 {
		printEmptyArgument(`Verification id`)
		return rcCommandLineError
	}

	rc := processCmdLineArguments(vcl, args[1:])
	if rc != rcOK {
		return rc
	}

	if vcl.BeQuiet {
		logger.SetLogLevel(logger.LogLevelWarning)
	}

	return doVerification(vcl.SignaturesFileName, verificationId)
}

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
		logger.PrintErrorf(mainMsgBase+2, `Error getting data from command line: %v`, err)
		return rcProcessError
	}

	return rcOK
}
