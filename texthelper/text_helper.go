package texthelper

// GetCountEnding returns the correct ending string for a number of items.
func GetCountEnding(n int) string {
	if n != 1 {
		return "s"
	} else {
		return ""
	}
}
