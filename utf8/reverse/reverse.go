//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var b strings.Builder

	replacementChar := '\uFFFD'

	for i := len(input) - 1; i >= 0; {
		r := replacementChar

		j := i + 1

		//Ищем руну среди 1-4 байтов
		for k := i; r == replacementChar && j-k <= 4; k-- {
			r, _ = utf8.DecodeRuneInString(input[k:j])

			if r != replacementChar {
				b.WriteRune(r)
				i = k - 1
			}
		}

		//Если не нашли руну среди 4 байтов, то пишем 4 replacementChar
		if r == replacementChar {
			for k := 0; k < 4; k++ {
				b.WriteRune(replacementChar)
			}

			i -= 4
		}
	}

	return b.String()
}
