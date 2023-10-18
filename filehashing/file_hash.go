package filehashing

import (
	"filesigner/contexthasher"
	"filesigner/filehelper"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
)

// ******** Public types ********

type FileHasher struct {
	hasher hash.Hash
}

// ******** Creation functions ********

// NewHasher Create new file hasher structure.
func NewHasher(contextId string) (*FileHasher, error) {
	hasher := contexthasher.NewContextHasher(sha3.New512(), contextId)

	return &FileHasher{hasher}, nil
}

// ******** Public functions ********

// HashFile calculates the hash value for one file.
func (fh *FileHasher) HashFile(filePath string) ([]byte, error) {
	hasher := fh.hasher

	err := hashFileContent(hasher, filePath)
	if err != nil {
		return nil, err
	}

	return fh.hasher.Sum(nil), nil
}

// ******** Private functions ********

// hashFileContent writes the content of a file to a hasher.
func hashFileContent(hasher hash.Hash, filePath string) error {
	var err error
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return err
	}
	defer filehelper.CloseFile(f)

	if _, err = io.Copy(hasher, f); err != nil {
		return err
	}

	return nil
}
