package utils

// StringInSlice checks if a string is present in a slice of strings.
// It returns true if the string is found, and false otherwise.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
