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

package cmdline

import (
	"bufio"
	"filesigner/filehelper"
	"filesigner/flaglist"
	"filesigner/set"
	"filesigner/signaturehandler"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ******** Private constants ********

const readFromStdInArg = "-"

// ******** Public functions ********

// FilesToProcess searches for the file names that match the command line options
func FilesToProcess(args []string, signatureFileName string) ([]string, signaturehandler.SignatureType, error) {
	signCmd := flag.NewFlagSet("sign", flag.ContinueOnError)

	var signatureTypeText string
	signCmd.StringVar(&signatureTypeText, "algorithm", "ed25519", "Signature algorithm (either 'ed25519' or 'ecdsap521')")
	signCmd.StringVar(&signatureTypeText, "a", "ed25519", "Short form of 'algorithm'")

	var signaturesFileName string
	signCmd.StringVar(&signaturesFileName, "signatures-file", signatureFileName, "Name of the file that receives the signatures")
	signCmd.StringVar(&signaturesFileName, "s", signatureFileName, "Short form of 'signatures-file'")

	var fromFileName string
	signCmd.StringVar(&fromFileName, "from-file", "", "Name of the file that contains a list of files to sign")
	signCmd.StringVar(&fromFileName, "f", "", "Short form of 'from-file'")

	var beQuiet bool
	signCmd.BoolVar(&beQuiet, "quiet", false, "Only write output if something goes wrong")
	signCmd.BoolVar(&beQuiet, "q", false, "Short form of 'quiet'")

	var doRecursion bool
	signCmd.BoolVar(&doRecursion, "recurse", false, "Only write output if something goes wrong")
	signCmd.BoolVar(&doRecursion, "r", false, "Short form of 'recurse'")

	var readStdIn bool
	signCmd.BoolVar(&readStdIn, "stdin", false, "Read list of files from stdin")
	signCmd.BoolVar(&readStdIn, "n", false, "Short form of 'stdin'")

	excludeFileList := flaglist.NewFileSystemFlagList()
	signCmd.Var(excludeFileList, "exclude-file", "Name of file to exclude from signing (may contain wild-cards).")
	signCmd.Var(excludeFileList, "xf", "Short for 'exclude-file'")

	includeFileList := flaglist.NewFileSystemFlagList()
	signCmd.Var(includeFileList, "include-file", "Name of file to include in signing may contain wild-cards)")
	signCmd.Var(includeFileList, "if", "Short for 'include-file'")

	excludeDirList := flaglist.NewFileSystemFlagList()
	signCmd.Var(excludeDirList, "exclude-dir", "Name of directory to exclude from signing (may contain wild-cards).")
	signCmd.Var(excludeDirList, "xd", "Short for 'exclude-dir'")

	includeDirList := flaglist.NewFileSystemFlagList()
	signCmd.Var(includeDirList, "include-dir", "Name of directory to include in signing may contain wild-cards)")
	signCmd.Var(includeDirList, "id", "Short for 'include-dir'")

	var signatureType signaturehandler.SignatureType

	// 1. Parse command line.
	err := signCmd.Parse(args)
	if err != nil {
		return nil, signatureType, err
	}

	// 2. Get signature type.
	signatureType, err = convertSignatureType(strings.ToLower(signatureTypeText))
	if err != nil {
		return nil, signatureType, err
	}

	// 3. Read file names from command line, StdIn and options.
	var fileSpecs []string
	fileSpecs, err = getFileSpecsFromCmdLine(signCmd.Args(), fromFileName, readStdIn)
	if err != nil {
		return nil, signatureType, err
	}

	// 4. Move any command line wild cards to the includeFileList.
	fileSpecs = moveWildCardFileSpecs(fileSpecs, includeFileList)

	// 5. Convert file specs to absolute path names.
	fileSpecs, err = makeAbsFileSpecs(fileSpecs)
	if err != nil {
		return nil, signatureType, err
	}

	// 6. Get the real path names for the file specifications.
	var filePaths *set.FileSystemStringSet
	filePaths, err = getRealFilePathsFromSpecs(fileSpecs, excludeDirList.Elements(), excludeFileList.Elements())
	if err != nil {
		return nil, signatureType, err
	}

	// 7. If no files are specified, or any include "include" is specified, scan the current directory.
	var scanPaths *set.FileSystemStringSet
	if filePaths.Len() == 0 || includeFileList.Len() != 0 || includeDirList.Len() != 0 {
		scanPaths, err = filehelper.ScanDir(includeFileList, excludeFileList, includeDirList, excludeDirList, doRecursion)
	} else {
		scanPaths = set.NewFileSystemStringSet()
	}

	// 8. Combine the two file lists and return.
	return filePaths.Union(scanPaths).Elements(), signatureType, nil
}

// ******** Private functions ********

// convertSignatureType converts the signature type text into a SignatureType value.
func convertSignatureType(signatureTypeText string) (signaturehandler.SignatureType, error) {
	switch signatureTypeText {
	case "ed25519":
		return signaturehandler.SignatureTypeEd25519, nil

	case "ecdsap521":
		return signaturehandler.SignatureTypeEcDsaP521, nil

	default:
		return signaturehandler.SignatureTypeInvalid, fmt.Errorf("Invalid signature type: '%s'", signatureTypeText)
	}
}

// moveWildCardFileSpecs moves wild card file specifications to the includeFileList
func moveWildCardFileSpecs(fileSpecs []string, includeFileList *flaglist.FileSystemFlagList) []string {
	resultList := make([]string, 0, len(fileSpecs))
	for _, fileSpec := range fileSpecs {
		if strings.ContainsAny(fileSpec, `*?`) {
			_ = includeFileList.Set(fileSpec)
		} else {
			resultList = append(resultList, fileSpec)
		}
	}

	return resultList
}

// getFileSpecsFromCmdLine gathers all file specifications from the command line.
func getFileSpecsFromCmdLine(args []string, fromFileName string, readStdIn bool) ([]string, error) {
	var err error
	var fileSpecs []string

	// 1. See if there is a file that contains file names.
	if len(fromFileName) != 0 {
		fileSpecs, err = addFileSpecsFromFileName(fromFileName, fileSpecs)
	}

	// 2. Add file names from StdIn and the command line.
	fileSpecs = addFileSpecsFromCmdLineAndStdIn(readStdIn, args, fileSpecs)
	return fileSpecs, err
}

// getRealFilePathsFromSpecs returns all file paths that match the supplied file specifications.
func getRealFilePathsFromSpecs(fileSpecs []string, excludeDirList []string, excludeFileList []string) (*set.FileSystemStringSet, error) {
	filePaths := set.NewFileSystemStringSetWithLength(len(fileSpecs))

	thisDirPath, err := makeThisDirPath()
	if err != nil {
		return filePaths, err
	}

	var selectedFilePaths []string
	for _, fileSpec := range fileSpecs {
		selectedFilePaths, err = filehelper.PathGlob(fileSpec, excludeDirList, excludeFileList)
		if err != nil {
			return nil, err
		}

		if len(selectedFilePaths) == 0 {
			return nil, fmt.Errorf("No files found for specification '%s'", fileSpec)
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
	return "", fmt.Errorf("file path '%s' is not inside current directory '%s'", filePath, thisDirPath)
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
