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
//    2024-02-01: V1.0.0: Created.
//    2024-02-07: V2.0.0: Make an object.
//    2024-04-05: V2.0.1: Make Stdout the output destination for usage messages.
//

package cmdline

import (
	"bufio"
	"errors"
	"filesigner/filehelper"
	"filesigner/flaglist"
	"filesigner/set"
	"filesigner/signaturehandler"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"strings"
)

// ******** Private constants ********

// readFromStdInArg is the character after the end of all options that means "read from stdin".
const readFromStdInArg = `-`

// defaultSignatureAlgorithm is the name of the default signature algorithm.
const defaultSignatureAlgorithm = `ed25519`

// Constants for error messages re. include and exclude lists.

const excludeType = `ex`
const includeType = `in`
const fileObject = `file`
const directoryObject = `directory`

// ******** Public types ********

// SignCommandLine is the object that contains all the data
// to interpret a "sign" command line.
type SignCommandLine struct {
	// Public elements
	FileList           []string
	SignaturesFileName string
	SignatureType      signaturehandler.SignatureType

	// Private elements
	fs                *pflag.FlagSet
	signatureTypeText string
	prefix            string
	fromFileName      string
	beQuiet           bool
	doRecursion       bool
	readStdIn         bool
	excludeFileList   *flaglist.FileSystemFlagList
	excludeDirList    *flaglist.FileSystemFlagList
	includeFileList   *flaglist.FileSystemFlagList
	includeDirList    *flaglist.FileSystemFlagList
}

// ******** Public functions ********

// NewSignCommandLine sets up the flag parser for the "sign" command.
func NewSignCommandLine() *SignCommandLine {
	signCmd := pflag.NewFlagSet(`sign`, pflag.ContinueOnError)

	signCmd.SetOutput(os.Stdout)

	result := &SignCommandLine{fs: signCmd}

	signCmd.StringVarP(&result.signatureTypeText, `algorithm`, `a`, defaultSignatureAlgorithm, `Signature algorithm (either 'ed25519' or 'ecdsap521')`)

	signCmd.StringVarP(&result.prefix, `name`, `m`, defaultSignaturesFileNamePrefix, `Prefix of the signatures file name`)

	signCmd.StringVarP(&result.fromFileName, `from-file`, `f`, ``, `Name of a file that contains a list of files to sign`)

	signCmd.BoolVarP(&result.doRecursion, `recurse`, `r`, false, `Search this directory and all subdirectories`)

	signCmd.BoolVarP(&result.readStdIn, `stdin`, `s`, false, `Read list of files from stdin`)

	result.excludeFileList = flaglist.NewFileSystemFlagList()
	signCmd.VarP(result.excludeFileList, `exclude-file`, `x`, `Name of file to exclude from signing (may contain wildcards).`)

	result.includeFileList = flaglist.NewFileSystemFlagList()
	signCmd.VarP(result.includeFileList, `include-file`, `i`, `Name of file to include in signing (may contain wildcards)`)

	result.excludeDirList = flaglist.NewFileSystemFlagList()
	signCmd.VarP(result.excludeDirList, `exclude-dir`, `X`, `Name of directory to exclude from signing (may contain wildcards).`)

	result.includeDirList = flaglist.NewFileSystemFlagList()
	signCmd.VarP(result.includeDirList, `include-dir`, `I`, `Name of directory to include in signing may contain wildcards)`)

	signCmd.SortFlags = true

	return result
}

// Parse parses the command line according to the flag rules.
func (cl *SignCommandLine) Parse(args []string) (error, bool) {
	err := cl.fs.Parse(args)
	if errors.Is(err, pflag.ErrHelp) {
		return nil, true
	}

	return err, false
}

// PrintUsage prints the usage information for the command.
func (cl *SignCommandLine) PrintUsage() {
	cl.fs.PrintDefaults()
}

// ExtractCommandData extracts the data that are needed for the command from the command line.
func (cl *SignCommandLine) ExtractCommandData() error {
	var err error

	// 1. Build signatures file name.
	cl.SignaturesFileName = cl.prefix + signaturesFileNameSuffix

	// 2. The signatures file must be written to the current directory.
	err = checkSignaturesFileName(cl.SignaturesFileName)
	if err != nil {
		return err
	}

	// 3. The signatures file must always be excluded.
	_ = cl.excludeFileList.Set(cl.SignaturesFileName)

	// 4. Get signature type.
	cl.SignatureType, err = convertSignatureType(strings.ToLower(cl.signatureTypeText))
	if err != nil {
		return err
	}

	// 5. Read file names from command line, StdIn and options.
	var fileSpecs []string
	fileSpecs, err = getFileSpecsFromCmdLine(cl.fs.Args(), cl.fromFileName, cl.readStdIn)
	if err != nil {
		return err
	}

	// 6. Move any command line wild cards to the includeFileList.
	fileSpecs = moveWildCardFileSpecs(fileSpecs, cl.includeFileList)

	// 7. Check for path separators in includes and excludes.
	err = checkExcludesIncludes(cl.excludeFileList.Elements(), cl.includeFileList.Elements(), cl.excludeDirList.Elements(), cl.includeDirList.Elements())
	if err != nil {
		return err
	}

	// 8. Convert file specs to absolute path names.
	fileSpecs, err = makeAbsFileSpecs(fileSpecs)
	if err != nil {
		return err
	}

	// 9. Get the real path names for the file specifications.
	var filePaths *set.Set[string]
	filePaths, err = getRealFilePathsFromSpecs(fileSpecs, cl.excludeDirList.Elements(), cl.excludeFileList.Elements())
	if err != nil {
		return err
	}

	// 10. If no files are specified, or any include "include" is specified, scan the current directory.
	var scanPaths *set.Set[string]
	if filePaths.Size() == 0 || cl.includeFileList.Size() != 0 || cl.includeDirList.Size() != 0 {
		scanPaths, err = filehelper.ScanDir(
			cl.includeFileList.Elements(),
			cl.excludeFileList.Elements(),
			cl.includeDirList.Elements(),
			cl.excludeDirList.Elements(),
			cl.doRecursion,
		)
		if err != nil {
			return err
		}
	} else {
		// scanPaths is an empty set if the directory is not scanned.
		scanPaths = set.New[string]()
	}

	// 11. Combine the two file lists and return.
	cl.FileList = filePaths.Union(scanPaths).Elements()

	return nil
}

// ******** Private functions ********

// convertSignatureType converts the signature type text into a SignatureType value.
func convertSignatureType(signatureTypeText string) (signaturehandler.SignatureType, error) {
	switch signatureTypeText {
	case `ed25519`:
		return signaturehandler.SignatureTypeEd25519, nil

	case `ecdsap521`:
		return signaturehandler.SignatureTypeEcDsaP521, nil

	default:
		return signaturehandler.SignatureTypeInvalid, fmt.Errorf(`Invalid signature type: '%s'`, signatureTypeText)
	}
}

// moveWildCardFileSpecs moves wild card file specifications to the includeFileList
func moveWildCardFileSpecs(fileSpecs []string, includeFileList *flaglist.FileSystemFlagList) []string {
	resultList := make([]string, 0, len(fileSpecs))
	for _, fileSpec := range fileSpecs {
		if strings.ContainsAny(fileSpec, wildCards) {
			_ = includeFileList.Set(fileSpec)
		} else {
			resultList = append(resultList, fileSpec)
		}
	}

	return resultList
}

// checkExcludesIncludes checks exclude and include list for path separators.
func checkExcludesIncludes(excludeFileList []string, includeFileList []string, excludeDirList []string, includeDirList []string) error {
	err := checkTypeList(excludeType, fileObject, excludeFileList)
	if err != nil {
		return err
	}

	err = checkTypeList(includeType, fileObject, includeFileList)
	if err != nil {
		return err
	}

	err = checkTypeList(excludeType, directoryObject, excludeDirList)
	if err != nil {
		return err
	}

	err = checkTypeList(includeType, directoryObject, includeDirList)
	if err != nil {
		return err
	}

	return nil
}

// checkTypeList checks if a pattern contains a path separator.
func checkTypeList(listType string, listObject string, excludeFileList []string) error {
	for _, pattern := range excludeFileList {
		if !filehelper.IsFileName(pattern) {
			return fmt.Errorf(`Pattern '%s' in %sclude %s option must be a file name pattern`, pattern, listType, listObject)
		}
	}

	return nil
}

// getFileSpecsFromCmdLine gathers all file specifications from the command line.
func getFileSpecsFromCmdLine(args []string, fromFileName string, readStdIn bool) ([]string, error) {
	var err error
	var fileSpecs []string

	// 1. See if there is a file that contains file names.
	if len(fromFileName) != 0 {
		fileSpecs, err = addFileSpecsFromFileName(fromFileName, fileSpecs)
		if err != nil {
			return nil, err
		}
	}

	// 2. Add file names from StdIn and the command line.
	fileSpecs = addFileSpecsFromCmdLineAndStdIn(readStdIn, args, fileSpecs)

	return fileSpecs, nil
}

// getRealFilePathsFromSpecs returns all file paths that match the supplied file specifications.
func getRealFilePathsFromSpecs(fileSpecs []string, excludeDirList []string, excludeFileList []string) (*set.Set[string], error) {
	filePaths := set.NewWithLength[string](len(fileSpecs))

	// err is an OS error
	thisDirPath, err := makeThisDirPath()
	if err != nil {
		return filePaths, err
	}

	var selectedFilePaths []string
	for _, fileSpec := range fileSpecs {
		// Err can be "bad pattern" or some OS error
		selectedFilePaths, err = filehelper.PathGlob(fileSpec, excludeDirList, excludeFileList)
		if err != nil {
			return nil, err
		}

		if len(selectedFilePaths) == 0 {
			return nil, fmt.Errorf(`No files found for specification '%s' (maybe excluded)`, fileSpec)
		}

		for _, selectedFilePath := range selectedFilePaths {
			selectedFilePath, err = removeThisDirPath(thisDirPath, selectedFilePath)
			if err != nil {
				return filePaths, err
			}

			filePaths.Add(selectedFilePath)
		}
	}

	return filePaths, nil
}

// removeThisDirPath removes the current directory path from a file path.
func removeThisDirPath(thisDirPath string, filePath string) (string, error) {
	if strings.HasPrefix(filePath, thisDirPath) {
		return filePath[len(thisDirPath):], nil
	}
	return "", fmt.Errorf(`File path '%s' is not inside current directory '%s'`, filePath, thisDirPath)
}

// makeThisDirPath builds the current directory path.
func makeThisDirPath() (string, error) {
	thisDirPath, err := os.Getwd()

	if err != nil {
		return ``, err
	}

	if !os.IsPathSeparator(thisDirPath[len(thisDirPath)-1]) {
		thisDirPath += string(os.PathSeparator)
	}

	filehelper.EnsureDriveLetterIsUpperCase(thisDirPath)

	return thisDirPath, nil
}

// makeAbsFileSpecs converts the supplied file specifications to absolute path specifications.
func makeAbsFileSpecs(fileSpecs []string) ([]string, error) {
	for i, fileSpec := range fileSpecs {
		// Make an absolute path.
		normalizedFileSpec, err := filepath.Abs(fileSpec)
		if err != nil {
			return nil, err
		}

		fileSpecs[i] = normalizedFileSpec
	}

	return fileSpecs, nil
}

// addFileSpecsFromCmdLineAndStdIn adds files from StdIn and from the command line.
func addFileSpecsFromCmdLineAndStdIn(readStdIn bool, args []string, fileSpecs []string) []string {
	// Are there any arguments?
	if len(args) != 0 {
		// If 1. argument on the command line is '-' set readStdIn.
		if args[0] == readFromStdInArg {
			readStdIn = true
			args = args[1:] // Set args to files specs remaining after '-', if any.
		}

		if len(args) != 0 {
			// Read files from command line
			fileSpecs = append(fileSpecs, args...)
		}
	}

	// Read file names from StdIn
	if readStdIn {
		fileSpecs = addFilesFromFile(os.Stdin, fileSpecs)
	}

	return fileSpecs
}

// addFileSpecsFromFileName reads the contents of the file with the supplied file name
// and adds them to the given fileLines slice. It returns the updated fileLines slice.
func addFileSpecsFromFileName(fromFileName string, fileSpecs []string) ([]string, error) {
	readFile, err := os.Open(fromFileName)

	if err != nil {
		return nil, err
	}

	defer filehelper.CloseFile(readFile)

	fileSpecs = addFilesFromFile(readFile, fileSpecs)

	return fileSpecs, nil
}

// addFilesFromFile reads the content of the given os.File and appends each line to the provided fileLines slice.
// It returns the updated fileLines slice.
func addFilesFromFile(fromFile *os.File, fileSpecs []string) []string {
	fileScanner := bufio.NewScanner(fromFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileSpecs = append(fileSpecs, fileScanner.Text())
	}

	return fileSpecs
}
