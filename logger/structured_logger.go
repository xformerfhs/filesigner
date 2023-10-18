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
