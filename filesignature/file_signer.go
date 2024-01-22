package filesignature

import (
	"filesigner/base32encoding"
	"filesigner/filehasher"
	"filesigner/hashsignature"
	"filesigner/maphelper"
	"fmt"
	"path/filepath"
)

// ******** Public functions ********

// SignFileHashes creates signatures for file hashes.
func SignFileHashes(hashSigner hashsignature.HashSigner,
	hashResultList map[string]*filehasher.HashResult) (map[string]string, []string, error) {
	filePathList := maphelper.SortedKeys(hashResultList)

	return makeHashSignatures(hashSigner, filePathList, hashResultList)
}

// ******** Private functions ********

func makeHashSignatures(hashSigner hashsignature.HashSigner,
	filePathList []string,
	hashResultList map[string]*filehasher.HashResult) (map[string]string, []string, error) {
	var err error
	signatures := make(map[string]string, len(filePathList))

	var signature []byte
	for _, filePath := range filePathList {
		signature, err = hashSigner.SignHash(hashResultList[filePath].HashValue)
		if err != nil {
			return nil, nil, fmt.Errorf("Could not sign hash of file '%s': %w", filePath, err)
		}

		signatures[filepath.ToSlash(filePath)] = base32encoding.EncodeToString(signature)
	}

	return signatures, filePathList, nil
}
