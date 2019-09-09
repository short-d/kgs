package unique

// Characters  returns true if chars contains only unique characters
func Characters(chars []byte) bool {
	seen := make(map[byte]bool)
	for _, char := range chars {
		if _, ok := seen[char]; ok {
			return false
		}
		seen[char] = true
	}
	return true
}
