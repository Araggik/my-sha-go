//go:build !solution

package reverse

import (
	"strings"
)

func Reverse(input string) string {
	var b strings.Builder

	runes := []rune{}

	for _, r := range input {
		runes = append(runes, r)
	}

	for i := len(runes) - 1; i >= 0; i-- {
		b.WriteRune(runes[i])
	}

	return b.String()
}
