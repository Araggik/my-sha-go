//go:build !solution

package varfmt

import (
	"fmt"
	"strings"
	"strconv"
)

func Sprintf(format string, args ...interface{}) string {
	var b strings.Builder

	braceChar1 := '{'
	braceChar2 := '}'

	isArgNumber := false

	var argNumberStr strings.Builder

	braceCount := -1

	for _, r := range format {
		switch  {
		  case r == braceChar1:
			isArgNumber = true

			braceCount++
		  case r == braceChar2 && isArgNumber:
			isArgNumber = false

			if argNumberStr.Len() == 0 {
				b.WriteString(fmt.Sprint(args[braceCount]))
			} else {
				argNumber, err := strconv.Atoi(argNumberStr.String())

				if err != nil {
					panic("Некорректный индекс аргумента")
				}
	
				b.WriteString(fmt.Sprint(args[argNumber]))
	
				argNumberStr.Reset()
			}
		  case isArgNumber:
			argNumberStr.WriteRune(r)
		  default:
			b.WriteRune(r)
		}
	}

	return b.String()
}
