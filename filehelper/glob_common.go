package filehelper

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
