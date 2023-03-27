package utils

// InStrings check the []string contains the given element
func InStrings(ss []string, val string) bool {
	for _, ele := range ss {
		if ele == val {
			return true
		}
	}
	return false
}
