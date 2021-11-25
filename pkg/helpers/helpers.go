package helpers

// IsInTypeArray checks if an element is present in an array.
func IsInTypeArray(element int, array []int) bool {
	for _, x := range array {
		if x == element {
			return true
		}
	}
	return false
}
