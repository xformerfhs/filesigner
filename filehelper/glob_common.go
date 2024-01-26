package filehelper

// ******** Public functions ********

// SensibleGlobFilesOnly globs a pattern with Windows API calls as this is the only correct
// way to handle globbing on the case-insensitive Windows file system.
// This call returns only files, no directories.
func SensibleGlobFilesOnly(pattern string) ([]string, error) {
	return sensibleGlobWithSwitch(pattern, false, true)
}

// SensibleGlobDirsOnly globs a pattern with Windows API calls as this is the only correct
// way to handle globbing on the case-insensitive Windows file system.
// This call returns only directories, no files.
func SensibleGlobDirsOnly(pattern string) ([]string, error) {
	return sensibleGlobWithSwitch(pattern, true, false)
}

// SensibleGlob globs a pattern with Windows API calls as this is the only correct
// way to handle globbing on the case-insensitive Windows file system.
func SensibleGlob(pattern string) ([]string, error) {
	return sensibleGlobWithSwitch(pattern, true, true)
}

// ******** Private functions ********

// ensureNoTrailingSeparator ensures that the pattern does not end with a trailing separator
func ensureNoTrailingSeparator(pattern string) string {
	var pos int
	for pos = len(pattern) - 1; pos >= 0; {
		b := pattern[pos]
		if (b == byte('\\')) || (b == byte('/')) {
			pos--
		} else {
			break
		}
	}

	return pattern[:pos+1]
}
