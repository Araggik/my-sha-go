//go:build !solution

package speller

import "strings"

var placeNumberMap = map[int]string{
	1: "billion",
	2: "million",
	3: "thousand",
}

var twentyMap = map[int]string{
	0: "zero",
	1: "one",
	2: "two",
	3: "three",
	4: "four",
	5: "five",
	6: "six",
	7: "seven",
	8: "eight",
	9: "nine",
	10: "ten",
	11: "eleven",
	12: "twelve",
	13: "thirteen",
	14: "fourteen",
	15: "fifteen",
	16: "sixteen",
	17: "seventeen",
	18: "eighteen",
	19: "nineteen",
}

var dozenMap = map[int]string{
	2: "twenty",
	3: "thirty",
	4: "forty",
	5: "fifty",
	6: "sixty",
	7: "seventy",
	8: "eighty",
	9: "ninety",
}

func receiveStrForTreeDigit(treeDigitNumber int) string{
	var b strings.Builder

	isHundred := treeDigitNumber > 100

	if isHundred {
		hundred := treeDigitNumber / 100

		b.WriteString(twentyMap[hundred])
		b.WriteRune(' ')
		b.WriteString("hundred")
	}

	twoDigit := treeDigitNumber % 100

	if twoDigit > 0{
		if isHundred {
			b.WriteRune(' ')
		}

		if twoDigit > 19 {
			digit := twoDigit % 10

			dozen := twoDigit / 10

			if digit == 0 {
				b.WriteString(dozenMap[dozen])
			} else {
				b.WriteString(dozenMap[dozen])
				b.WriteRune('-')
				b.WriteString(twentyMap[digit])
			}
		} else {
			b.WriteString(twentyMap[twoDigit])
		}
	}

	return b.String()
}

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}

	var b strings.Builder

	var x int64

	if n < 0 {
		x = -n

		b.WriteString("minus ")
	} else {
		x = n
	}

	var deriveNumber int64 = 1000000000

	for i := 1; i <= 4; i++ {
		treeDigitNumber := int(n / deriveNumber)

		if treeDigitNumber > 0 {
			s := receiveStrForTreeDigit(treeDigitNumber)

			b.WriteRune(' ')
			b.WriteString(s)

			if i != 4 {
				b.WriteRune(' ')
				b.WriteString(placeNumberMap[i])
	
				x -= int64(treeDigitNumber) * deriveNumber
			}	
		}

		deriveNumber /= 1000	
	}

	return b.String()
}
