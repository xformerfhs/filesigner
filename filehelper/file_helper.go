package filehelper

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ******** Public functions ********

// CloseFile closes a file and prints an error message if closing failed.
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("error closing file '%s': %v", file.Name(), err)
	}
}

// GetRealBaseName gets the base name of a file without the extension.
func GetRealBaseName(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}

// FileSize returns the size of the named file.
func FileSize(filePath string) (int64, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
