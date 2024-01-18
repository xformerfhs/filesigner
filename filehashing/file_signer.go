package filehashing

import (
	"filesigner/base32encoding"
	"filesigner/hashsigner"
	"filesigner/maphelper"
	"fmt"
)

// ******** Public functions ********

// SignFileHashes creates signatures for file hashes.
func SignFileHashes(hashSigner hashsigner.HashSigner,
	hashResultList map[string]*HashResult) (map[string]string, error) {
	filePathList := maphelper.SortedKeys(hashResultList)

	return makeHashSignatures(hashSigner, filePathList, hashResultList)
}

// ******** Private functions ********

func makeHashSignatures(hashSigner hashsigner.HashSigner,
	filePathList []string,
	hashResultList map[string]*HashResult) (map[string]string, error) {
	var err error
	signatures := make(map[string]string, len(filePathList))

	var signature []byte
	for _, filePath := range filePathList {
		signature, err = hashSigner.SignHash(hashResultList[filePath].HashValue)
		if err != nil {
			return nil, fmt.Errorf("Could not sign hash of file '%s': %w", filePath, err)
		}

		signatures[filePath] = base32encoding.EncodeToString(signature)
	}

	return signatures, nil
}
