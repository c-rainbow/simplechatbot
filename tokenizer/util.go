package tokenizer

// IsSubRune If runes is equal to subrunes from fromIndex.
func IsSubRune(runes []rune, subrunes []rune, fromIndex int) bool {
	if len(runes) < fromIndex+len(subrunes) {
		return false
	}
	// compare all indexes of runes starting from "fromIndex"
	for i, subRune := range subrunes {
		if runes[fromIndex+i] != subRune {
			return false
		}
	}
	return true
}
