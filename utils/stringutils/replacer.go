package stringutils

// IStringReplacer applies a set of replacements to a string.
type IStringReplacer interface {
	// Replace returns a copy of s with all replacements performed.
	Replace(s string) string
}
