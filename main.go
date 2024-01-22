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

// signatureFileName is the fixed file name of the signature file
const signatureFileName = "signatures.json"

// Program information

// myVersion contains the program version.
const myVersion = "0.50.0"

// myName contains the program name.
var myName string

// ******** Real main function ********

// mainWithReturnCode is the real main function with arguments and return code.
func mainWithReturnCode(args []string) int {
	myName = filehelper.GetRealBaseName(args[0])

	printVersion()

	argLen := len(args)
	if argLen < 2 {
		return printUsageError(11, "Not enough arguments")
	}

	command := strings.ToLower(args[1])

	switch command {
	case commandHelp:
		printUsageText()
		return rcOK

	case commandSign:
		if argLen < 3 {
			return printMissingSignParameters("Context id")
		}

		fileList, signatureType, err := cmdline.FilesToProcess(args[3:], signatureFileName)
		if err != nil {
			return printUsageErrorf(12, "Error processing file names: %v", err)
		}

		return doSigning(signatureType, args[2], fileList)

	case commandVerify:
		if argLen < 3 {
			return printUsageError(12, "Context id missing")
		}

		if argLen > 3 {
			return printUsageError(13, "There must be no files specified for verification")
		}

		return doVerification(args[2])

	default:
		return printUsageErrorf(14, "Unknown command: '%s'", command)
	}
}

// ******** Private functions ********

// printMissingSignParameters prints an error message for missing sign parameters.
func printMissingSignParameters(parameters string) int {
	return printUsageErrorf(17, "%s of files to sign missing", parameters)
}

// printVersion prints the program version information.
func printVersion() {
	logger.PrintInfof(18, "%s V%s (%s, %d cpus)",
		myName,
		myVersion,
		runtime.Version(),
		runtime.NumCPU())
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
	_, _ = fmt.Printf("\nUsage:\n\n   %s sign {contextId} [-type {type}] [-if|-include-file {mask}] [-xf|-exclude-file {mask}] [-id|-include-dir {mask}] [-xd|-exclude-dir {mask}] [-no-subdirs]\n", myName)
	_, _ = fmt.Printf("      sign: Sign files and write signatures into file '%s'\n", signatureFileName)
	_, _ = fmt.Println("           contextId:    Arbitrary string used as a domain separator")
	_, _ = fmt.Println("           type:         Signature type (optional, 'ed25519' or 'ecdsap521', default is 'ed25519')")
	_, _ = fmt.Println("           include-file: File to include (optional, may contain wildcards, one per option)")
	_, _ = fmt.Println("              if:        Short for 'include-file'")
	_, _ = fmt.Println("           exclude-file: File to exclude (optional, may contain wildcards, one per option)")
	_, _ = fmt.Println("              xf:        Short for 'exclude-file'")
	_, _ = fmt.Println("           include-dir:  Directory to include (optional, may contain wildcards, one per option)")
	_, _ = fmt.Println("              id:        Short for 'include-dir'")
	_, _ = fmt.Println("           exclude-dir:  Directory to exclude (optional, may contain wildcards, one per option)")
	_, _ = fmt.Println("              xd:        Short for 'exclude-dir'")
	_, _ = fmt.Println("           no-subdirs:   Do not descend into subdirectories (optional)")
	if runtime.GOOS != "windows" {
		_, _ = fmt.Println("      Masks with wildcards need to be enclosed in quotes (') or double quotes (\")")
	}
	_, _ = fmt.Println("      Specifying an 'include' option implies that all others are excluded")
	_, _ = fmt.Printf("\n   %s verify {contextId}\n", myName)
	_, _ = fmt.Printf("      verify: Verify files with signatures in file '%s'\n", signatureFileName)
	_, _ = fmt.Println("         contextId: Arbitrary string used as a domain separator")
	_, _ = fmt.Println()
}
