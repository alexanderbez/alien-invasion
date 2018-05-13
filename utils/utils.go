package utils

// RemoveStrFromSlice attempts to remove a given string from a slice of
// strings. The given slice of strings is modified if the desired match string
// is removed.
func RemoveStrFromSlice(s []string, strMatch string) []string {
	for i, str := range s {
		if str == strMatch {
			s = append(s[:i], s[i+1:]...)
			break
		}
	}

	return s
}
