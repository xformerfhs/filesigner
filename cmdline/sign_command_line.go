package cmdline

import (
	"filesigner/filehelper"
	"filesigner/signaturehandler"
	"filesigner/typelist"
	"flag"
	"fmt"
	"strings"
)

// FilesToProcess searches for the file names that match the command line options
func FilesToProcess(args []string, signatureFileName string) ([]string, signaturehandler.SignatureType, error) {
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)

	var signatureTypeText string
	signCmd.StringVar(&signatureTypeText, "type", "ed25519", "Signature type (either 'ed25519' or 'ecdsap521')")

	var noSubDirs bool
	signCmd.BoolVar(&noSubDirs, "no-subdirs", false, "Do not process subdirectories")

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

	var signatureType signaturehandler.SignatureType
	err := signCmd.Parse(args)
	if err != nil {
		return nil, signatureType, err
	}

	if len(signCmd.Args()) != 0 {
		return nil, signatureType, fmt.Errorf("Extraneous parameters on command line: %v\nUse only one mask per option", signCmd.Args())
	}

	*excludeFileList = append(*excludeFileList, signatureFileName)

	signatureType, err = convertSignatureType(strings.ToLower(signatureTypeText))
	if err != nil {
		return nil, signatureType, err
	}

	var resultList []string
	resultList, err = filehelper.ScanDir(includeFileList,
		excludeFileList,
		includeDirList,
		excludeDirList,
		noSubDirs)

	if err != nil {
		return nil, signatureType, err
	}

	return resultList, signatureType, nil
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
