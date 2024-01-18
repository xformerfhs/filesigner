package filesignature

import (
	"filesigner/base32encoding"
	"filesigner/filehasher"
	"filesigner/hashsignature"
	"filesigner/maphelper"
	"fmt"
)

// ******** Public functions ********

// VerifyFileHashes verifies file hashes
func VerifyFileHashes(hashVerifier hashsignature.HashVerifier,
	fileSignatures map[string]string,
	fileHashList map[string]*filehasher.HashResult) ([]string, []error) {
	var err error

	successCollection := make([]string, 0, len(fileHashList))
	errCollection := make([]error, 0, len(fileHashList))

	filePathList := maphelper.SortedKeys(fileSignatures)

	var signatureString string
	var signatureValue []byte
	for _, filePath := range filePathList {
		fileHashResult, haveHashForFilePath := fileHashList[filePath]
		if haveHashForFilePath {
			signatureString = fileSignatures[filePath]
			signatureValue, err = base32encoding.DecodeFromString(signatureString)
			if err != nil {
				errCollection = append(errCollection, fmt.Errorf("Signature of file '%s' has invalid encoding: %w", filePath, err))
			} else {
				var ok bool
				ok, err = hashVerifier.VerifyHash(fileHashResult.HashValue, signatureValue)
				if ok {
					successCollection = append(successCollection, filePath)
				} else {
					errCollection = append(errCollection, fmt.Errorf("File '%s' has been tampered with", filePath))
				}
			}
		}
	}

	return successCollection, errCollection
}
