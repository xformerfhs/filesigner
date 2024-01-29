package cmdline

import (
	"bufio"
	"filesigner/filehelper"
	"filesigner/signaturehandler"
	"filesigner/stringhelper"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const readFromStdInArg = "-"

// FilesToProcess searches for the file names that match the command line options
func FilesToProcess(args []string, signatureFileName string) ([]string, signaturehandler.SignatureType, error) {
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)

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
	/*
		excludeFileList := typelist.NewFlagTypeList()
		signCmd.Var(excludeFileList, "exclude-file", "Name of file to exclude from signing (may contain wild-cards).")
		signCmd.Var(excludeFileList, "xf", "Short for 'exclude-file'")

		includeFileList := typelist.NewFlagTypeList()
		signCmd.Var(includeFileList, "include-file", "Name of file to include in signing may contain wild-cards)")
		signCmd.Var(includeFileList, "if", "Short for 'include-file'")

		excludeDirList := typelist.NewFlagTypeList()
		signCmd.Var(excludeDirList, "exclude-dir", "Name of directory to exclude from signing.")
		signCmd.Var(excludeDirList, "xd", "Short for 'exclude-dir'")

		includeDirList := typelist.NewFlagTypeList()
		signCmd.Var(includeDirList, "include-dir", "Name of directory to include in signing")
		signCmd.Var(includeDirList, "id", "Short for 'include-dir'")
	*/

	var signatureType signaturehandler.SignatureType
	err := signCmd.Parse(args)
	if err != nil {
		return nil, signatureType, err
	}

	signatureType, err = convertSignatureType(strings.ToLower(signatureTypeText))
	if err != nil {
		return nil, signatureType, err
	}

	//	*excludeFileList = append(*excludeFileList, signatureFileName)

	resultList := make([]string, 0, 100)

	if len(fromFileName) != 0 {
		resultList, err = addFilesFromFileName(fromFileName, resultList)
	}

	if signCmd.NArg() != 0 {
		resultList = addFilesFromCmdLineOrStdIn(signCmd, resultList)
	}

	/*
			resultList, err = filehelper.ScanDir(includeFileList,
				excludeFileList,
				includeDirList,
				excludeDirList,
				noSubDirs)

			if err != nil {
				return nil, signatureType, err
			}

		return resultList, signatureType, nil
	*/
	resultList, err = checkAndNormalizeFilePaths(resultList)

	if err != nil {
		return nil, signatureType, err
	}

	resultList, err = convertFileSpecToDirNames(resultList)
	return resultList, signatureType, nil
}

func convertFileSpecToDirNames(cmdLinePaths []string) ([]string, error) {
	result := make([]string, 0, len(cmdLinePaths))
	var globList []string
	var err error
	for _, path := range cmdLinePaths {
		globList, err = filehelper.PathGlob(path)
		if err != nil {
			return nil, err
		}

		result = append(result, globList...)
	}

	return result, nil
}

// addFilesFromCmdLineOrStdIn adds files from the command line or from StdIn
func addFilesFromCmdLineOrStdIn(signCmd *flag.FlagSet, resultList []string) []string {
	remArgs := signCmd.Args()

	// Add files from the command line, except when the first argument is "-"
	if remArgs[0] != readFromStdInArg {
		resultList = append(resultList, remArgs...)
	} else {
		resultList = addFilesFromFile(os.Stdin, resultList)
	}

	return resultList
}

// checkAndNormalizeFilePaths checks the file paths if they are valid file paths and normalizes them.
func checkAndNormalizeFilePaths(filepathList []string) ([]string, error) {
	thisDirPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	thisDirPathLen := getPathLen(thisDirPath)

	var normalizedFilePath string

	resultList := make([]string, 0, len(filepathList))
	for _, filePath := range filepathList {
		// 1. Make file path a normalized absolute path
		normalizedFilePath, err = normalizeFilePath(filePath)
		if err != nil {
			return nil, err
		}

		// 2. Convert to relative path, or complain if the specified absolute path is not in the current directory
		if strings.HasPrefix(normalizedFilePath, thisDirPath) {
			resultList = append(resultList, normalizedFilePath[thisDirPathLen:])
		} else {
			return nil, fmt.Errorf("Absolute file path '%s' is not within the current directory", filePath)
		}
	}

	return resultList, nil
}

// normalizeFilePath converts a file path to an absolute file path.
func normalizeFilePath(filePath string) (string, error) {
	// 1. Convert the specified path to an absolute path
	normalizedFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	// 2. If this is Windows convert the drive letter to upper case
	//    This is a hack to avoid creating a new string just because
	//    the first character is converted to upper case
	normalizedFilePathBytes := stringhelper.UnsafeStringBytes(normalizedFilePath)
	if runtime.GOOS == "windows" && normalizedFilePathBytes[1] == ':' {
		if normalizedFilePathBytes[0] > 'Z' {
			normalizedFilePathBytes[0] ^= 0x20
		}
	}

	return normalizedFilePath, nil
}

// getPathLen calculates the length of the path that is used to convert an absolute to a relative path.
func getPathLen(filePath string) int {
	filePathLen := len(filePath)
	vol := filepath.VolumeName(filePath)
	filePath = filePath[len(vol):]

	// Add one to length to account for the file path separator.
	// Except when this is a root directory.
	if filePathLen > 1 || (filePath != `/` && filePath != `\`) {
		filePathLen++
	}

	return filePathLen
}

// convertSignatureType converts the signature type text into a byte.
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

func addFilesFromFileName(fromFileName string, fileLines []string) ([]string, error) {
	readFile, err := os.Open(fromFileName)

	if err != nil {
		return nil, err
	}

	defer filehelper.CloseFile(readFile)

	fileLines = addFilesFromFile(readFile, fileLines)

	return fileLines, nil
}

func addFilesFromFile(readFile *os.File, fileLines []string) []string {
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	return fileLines
}
