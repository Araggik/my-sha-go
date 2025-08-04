//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
)

func CollapseSpaces(input string) string {
	var b strings.Builder

	spaceChar := ' '

	isSpaceChars := false

	for _, r := range input {
		if unicode.IsSpace(r) {
			isSpaceChars = true
		} else {
			if isSpaceChars {
				isSpaceChars = false

				b.WriteRune(spaceChar)
			}

			b.WriteRune(r)
		}
	}

	return b.String()
}
