package editcorrector

import "strconv"

// unescapeGoString unescapes common Go string literal escape sequences.
// It handles \n, \r, \t, \", \\, and \uXXXX (Unicode).
// It's designed to correct strings that might have been double-escaped or
// incorrectly escaped when passed as string literals.
func UnescapeGoString(inputString string) string {
	// strconv.Unquote can unquote a Go string literal.
	// We need to add double quotes around the input to make it a valid string literal.
	quotedInput := `"` + inputString + `"`
	
	unquoted, err := strconv.Unquote(quotedInput)
	if err != nil {
		// If unquoting fails, it means it's not a valid Go string literal
		// or already unescaped. Return the original string.
		return inputString
	}
	return unquoted
}
