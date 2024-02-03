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

package logger

import (
	"fmt"
	"time"
)

// ******** Public types ********

type LogLevel byte

// ******** Public constants ********

const (
	LogLevelInfo LogLevel = iota
	LogLevelWarning
	LogLevelError
)

// ******** Private constants ********

const sevInfo byte = 'I'
const sevWarning byte = 'W'
const sevError byte = 'E'

// timeFormat is the time format for log messages.
const timeFormat = "2006-01-02 15:04:05 Z07:00"

// ******** Private variables ********

var logLevel LogLevel

// ******** Public functions ********

// SetLogLevel sets the log level.
func SetLogLevel(newLogLevel LogLevel) {
	if newLogLevel < LogLevelInfo {
		newLogLevel = LogLevelInfo
	} else {
		if newLogLevel > LogLevelError {
			newLogLevel = LogLevelError
		}
	}

	logLevel = newLogLevel
}

// -------- Text functions --------

// PrintInfo prints an information message.
func PrintInfo(msgNum byte, msgText string) {
	if logLevel <= LogLevelInfo {
		printLogLine(msgNum, sevInfo, msgText)
	}
}

// PrintWarning prints a warning message.
func PrintWarning(msgNum byte, msgText string) {
	if logLevel <= LogLevelWarning {
		printLogLine(msgNum, sevWarning, msgText)
	}
}

// PrintError prints an error message.
func PrintError(msgNum byte, msgText string) {
	if logLevel <= LogLevelError {
		printLogLine(msgNum, sevError, msgText)
	}
}

// -------- Format functions --------

// PrintInfof prints an information message with a format string.
func PrintInfof(msgNum byte, msgFormat string, args ...any) {
	PrintInfo(msgNum, fmt.Sprintf(msgFormat, args...))
}

// PrintWarningf prints a warning message with a format string.
func PrintWarningf(msgNum byte, msgFormat string, args ...any) {
	PrintWarning(msgNum, fmt.Sprintf(msgFormat, args...))
}

// PrintErrorf prints an error message with a format string.
func PrintErrorf(msgNum byte, msgFormat string, args ...any) {
	PrintError(msgNum, fmt.Sprintf(msgFormat, args...))
}

// ******** Private functions ********

// printLogLine prints the log line.
func printLogLine(msgNum byte, severity byte, msgText string) {
	fmt.Printf("%s  %d  %c  %s\n", time.Now().Format(timeFormat), msgNum, severity, msgText)
}
