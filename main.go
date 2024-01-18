package main

import (
	"filesigner/cmdline"
	"filesigner/filehelper"
	"filesigner/logger"
	"filesigner/signaturehandler"
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

// argumentSeparator separates parameters from the sign command
const argumentSeparator = "!"

// Program information

// myVersion contains the program version.
const myVersion = "0.14.0"

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
			return printMissingSignParameters("Context id, separator and list")
		}

		if argLen < 4 {
			return printMissingSignParameters("Separator and list")
		}

		fileListIndex, signatureTypeText, rc := checkSignatureTypeAndSeparator(args, argLen)
		if rc != rcOK {
			return rc
		}

		var signatureType signaturehandler.SignatureType
		signatureType, rc = convertSignatureType(signatureTypeText)
		if rc != rcOK {
			return rc
		}

		if argLen < fileListIndex+1 {
			return printMissingSignParameters("List")
		}

		args = append(args, string(cmdline.NegatePrefix)+signatureFileName) // Never include the signatures file
		return doSigning(signatureType, args[2], args[fileListIndex:])

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

// checkSignatureTypeAndSeparator checks for the presence of the separator and gets a signature type, if present.
func checkSignatureTypeAndSeparator(args []string, argLen int) (int, string, int) {
	fileListIndex := 4
	signatureTypeText := "ed25519"

	if args[3] != argumentSeparator {
		if argLen > 4 {
			signatureTypeText = strings.ToLower(args[3])

			if args[4] != argumentSeparator {
				return 0, "", printMissingSeparatorError()
			}

			fileListIndex++
		} else {
			return 0, "", printMissingSeparatorError()
		}
	}

	return fileListIndex, signatureTypeText, rcOK
}

// convertSignatureType converts the signature type text into a byte.
func convertSignatureType(signatureTypeText string) (signaturehandler.SignatureType, int) {
	switch signatureTypeText {
	case "ed25519":
		return signaturehandler.SignatureTypeEd25519, rcOK

	case "ecdsap521":
		return signaturehandler.SignatureTypeEcDsaP521, rcOK

	default:
		return signaturehandler.SignatureTypeInvalid, printUsageErrorf(16, "Invalid signature type: '%s'", signatureTypeText)
	}
}

// printMissingSeparatorError prints a message that the list separator is missing.
func printMissingSeparatorError() int {
	return printUsageError(15, "Required file list separator '!' missing")
}

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
	_, _ = fmt.Printf("\nUsage:\n\n   %s sign {contextId} [type] ! {fileList}\n", myName)
	_, _ = fmt.Printf("      sign: Sign files and write signatures into file '%s'\n", signatureFileName)
	_, _ = fmt.Println("         contextId: Arbitrary string used as a salt")
	_, _ = fmt.Println("         type: Signature type (optional, 'ed25519' or 'ecdsap521', default is 'ed25519')")
	_, _ = fmt.Println("         !:         Required separator before file list")
	_, _ = fmt.Println("         fileList:  Space-separated list of names of files to sign (wildcards are permitted)")
	_, _ = fmt.Println("                    Names prefixed by '-' are excluded from signatures")
	_, _ = fmt.Printf("\n   %s verify {contextId}\n", myName)
	_, _ = fmt.Printf("      verify: Verify files with signatures in file '%s'\n", signatureFileName)
	_, _ = fmt.Println("         contextId: Arbitrary string used as a salt")
	_, _ = fmt.Println()
}
